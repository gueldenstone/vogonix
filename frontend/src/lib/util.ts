function formatDurationFromNanosecondsInternal(nanoseconds: number): string[] {
  if (nanoseconds < 0) {
    throw new Error("Duration cannot be negative");
  }
  if (nanoseconds == 0) {
    return ['0s']
  }
  const units = [
    { label: 'y', seconds: 365 * 24 * 60 * 60 }, // years
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

  return parts
}


export function formatDurationFromNanoseconds(nanoseconds: number): string {
  return formatDurationFromNanosecondsInternal(nanoseconds).join(' ');
}

export function formatDuration(seconds: number): string {
  return formatDurationFromNanosecondsInternal(seconds * 1e9).join(' ')
}

export function timeAgo(inputTime: string): string {
  const past = new Date(inputTime);
  let diffInNanoseconds = Math.floor((Date.now().valueOf() - past.valueOf()) * 1e6);
  if (diffInNanoseconds < (1e9 * 60)) {
    return "just now"; // For times less than 1 minute ago
  }
  return formatDurationFromNanosecondsInternal(diffInNanoseconds)[0] + " ago";
}
