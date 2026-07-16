export function feelingLabel(feeling: string | null | undefined): string {
  const trimmedFeeling = feeling?.trim();
  return trimmedFeeling?.length ? trimmedFeeling : "No feeling recorded";
}
