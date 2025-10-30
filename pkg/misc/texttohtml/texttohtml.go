package texttohtml

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func ConvertAllHandlesToUrls(kernel interfaces.Kernel, handles interface{}) (map[string]string, error) {
	var handlesT []string
	switch v := handles.(type) {
	case []interface{}:
		for _, handle := range v {
			if handleStr, ok := handle.(string); ok {
				handlesT = append(handlesT, handleStr)
			} else {
				return nil, fmt.Errorf("invalid handle type")
			}
		}
	case []string:
		handlesT = v
	case string:
		handlesT = []string{v}
	default:
		return nil, fmt.Errorf("invalid handles type")
	}
	urls := make(map[string]string)
	for _, handle := range handlesT {
		if strings.HasPrefix(handle, "https://") || handle == "Public" || handle == "Followers" {
			urls[handle] = handle
		} else if strings.HasPrefix(handle, "@") {
			_, url, err := kernel.GetActorUrlByHandler(handle)
			if err != nil {
				return nil, fmt.Errorf("failed to convert handle %s: %v", handle, err)
			}
			urls[handle] = url
		}
	}
	return urls, nil
}

func TextToHtml(kernel interfaces.Kernel, domain string, input string) (string, []string, map[string]string) {
	escaped := html.EscapeString(input)
	hashtags := make([]string, 0)
	mentions := make(map[string]string)

	// Mentions: start of line or space
	pubMentionRe := regexp.MustCompile(`(^|\s)@([\w.]+)(?:@([\w.]+))?`)
	escaped = pubMentionRe.ReplaceAllStringFunc(escaped, func(m string) string {
		prefix := m[:1]
		handle := m[len(prefix):]
		if !strings.Contains(handle, "@") {
			handle = handle + "@" + domain
		}
		if !strings.HasPrefix(handle, "@") {
			handle = "@" + handle
		}
		if strings.HasSuffix(prefix, "@") {
			prefix = prefix[:len(prefix)-1]
		}
		actor, err := kernel.ActivityPubDB().GetActor("", handle)
		if err == nil {
			mentions[handle] = actor.Id
			return fmt.Sprintf("%s<a href=\"%s\" class=\"mention\">%s</a>", prefix, actor.Id, handle)
		} else {
			fmt.Println("ERROR ON TextToHtml: " + err.Error())
		}
		return m
	})

	// Hashtags
	hashtagRe := regexp.MustCompile(`(^|\s)#(\w+)`)
	escaped = hashtagRe.ReplaceAllStringFunc(escaped, func(m string) string {
		prefix := m[:1]
		tag := m[len(prefix)+1:]
		hashtags = append(hashtags, tag)
		return fmt.Sprintf("%s<a href=\"https://instance/tags/%s\" class=\"hashtag\">#%s</a>", prefix, tag, tag)
	})

	// URLs
	urlRe := regexp.MustCompile(`(^|\s)(https?://\S+)`)
	escaped = urlRe.ReplaceAllStringFunc(escaped, func(m string) string {
		prefix := m[:len(m)-len(m[strings.Index(m, "http"):])]
		url := m[len(prefix):]
		return fmt.Sprintf("%s<a href=\"%s\">%s</a>", prefix, url, url)
	})

	// Emoji
	emojiRe := regexp.MustCompile(`(^|\s):([\w_]+):`)
	escaped = emojiRe.ReplaceAllStringFunc(escaped, func(m string) string {
		prefix := m[:len(m)-len(":"+m[1:])]
		name := m[len(prefix)+1 : len(m)-1]
		return fmt.Sprintf("%s<span class=\"emoji\" role=\"img\" aria-label=\"%s\">:%s:</span>", prefix, name, name)
	})

	escaped = strings.ReplaceAll(escaped, "\n", "</p><p>")
	result := fmt.Sprintf(`<div class="e-content"><p>%s</p></div>`, escaped)
	if result == `<div class="e-content"><p></p></div>` {
		return "", hashtags, mentions
	}
	return result, hashtags, mentions
}
