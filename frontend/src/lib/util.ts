export function formatDuration(seconds: number): string {
  if (seconds < 0) {
    throw new Error("Duration cannot be negative");
  }

  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  return [
    hours > 0 ? `${hours}h` : null,
    minutes > 0 ? `${minutes}m` : null,
    `${secs}s`
  ]
    .filter(Boolean)
    .join(" ");
}
