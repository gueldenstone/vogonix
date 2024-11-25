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

export function formatDurationFromNanoseconds(nanoseconds: number): string {
  if (nanoseconds == 0) {
    return '0s'
  }
  const units = [
    { label: 'w', seconds: 7 * 24 * 60 * 60 }, // weeks
    { label: 'd', seconds: 24 * 60 * 60 },     // days
    { label: 'h', seconds: 60 * 60 },          // hours
    { label: 'm', seconds: 60 },               // minutes
    { label: 's', seconds: 1 },                // seconds
  ];

  // Convert nanoseconds to seconds
  let remainingSeconds = Math.floor(nanoseconds / 1e9);
  const parts: string[] = [];

  for (const unit of units) {
    const count = Math.floor(remainingSeconds / unit.seconds);
    if (count > 0) {
      parts.push(`${count}${unit.label}`);
      remainingSeconds %= unit.seconds;
    }
  }

  return parts.join(' ');
}
