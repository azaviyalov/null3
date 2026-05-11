const explicitEmojiPattern = /\p{Extended_Pictographic}/u;

const emojiRules: readonly { emoji: string; pattern: RegExp }[] = [
  {
    emoji: "🙂",
    pattern:
      /\b(happy|good|great|joy|glad|hopeful|grateful|better|content|fine|okay|ok|bright|proud)\b/i,
  },
  {
    emoji: "😌",
    pattern:
      /\b(calm|quiet|soft|gentle|rested|peace|peaceful|steady|balanced|relaxed|ease)\b/i,
  },
  {
    emoji: "😟",
    pattern: /\b(anx|worr|panic|nervous|overwhelm|restless|uneasy)\b/i,
  },
  {
    emoji: "😤",
    pattern: /\b(angry|mad|frustrat|annoy|irritat|tense|stress)\b/i,
  },
  {
    emoji: "😔",
    pattern:
      /\b(sad|low|down|blue|lonely|hurt|cry|heavy|tired|drained|exhausted|empty)\b/i,
  },
];

export function moodEmoji(feeling: string | null | undefined): string {
  const trimmedFeeling = feeling?.trim();
  if (!trimmedFeeling) {
    return "◌";
  }

  const matchedEmoji = trimmedFeeling.match(explicitEmojiPattern);
  if (matchedEmoji) {
    return matchedEmoji[0];
  }

  const matchingRule = emojiRules.find((rule) =>
    rule.pattern.test(trimmedFeeling),
  );
  return matchingRule?.emoji ?? "◌";
}

export function feelingLabel(feeling: string | null | undefined): string {
  const trimmedFeeling = feeling?.trim();
  return trimmedFeeling?.length ? trimmedFeeling : "Untitled feeling";
}
