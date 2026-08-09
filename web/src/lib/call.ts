export function validRoom(value: unknown): value is string {
  return typeof value === "string" && /^[A-Za-z0-9]{6}$/.test(value)
}

export function createRoomID(): string {
  const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
  let room = ""

  while (room.length < 6) {
    for (const byte of crypto.getRandomValues(new Uint8Array(8))) {
      room += alphabet[byte % alphabet.length]
      if (room.length === 6) break
    }
  }

  return room
}

export function cleanName(value: unknown): string {
  if (typeof value !== "string") return ""
  return value.replace(/[\u0000-\u001f\u007f]/g, "").trim().slice(0, 32)
}

export function detectDeviceType(): "mobile" | "computer" {
  if (navigator.userAgentData?.mobile) return "mobile"
  if (/Android|iPhone|iPad|iPod|IEMobile|Opera Mini/i.test(navigator.userAgent)) return "mobile"
  if (/Macintosh/i.test(navigator.userAgent) && navigator.maxTouchPoints > 1) return "mobile"
  return "computer"
}

export function videoConstraints(facingMode: VideoFacingModeEnum, exact = false): MediaTrackConstraints {
  return {
    width: { ideal: 1280 },
    height: { ideal: 720 },
    facingMode: exact ? { exact: facingMode } : { ideal: facingMode },
  }
}

export function friendlyMediaError(error: unknown): string {
  const issue = asError(error)
  if (issue.name === "NotAllowedError") return "Camera and microphone access is needed to join."
  if (issue.name === "NotFoundError") return "I couldn't find a camera or microphone on this device."
  if (issue.name === "NotReadableError" || issue.name === "AbortError") {
    return "The camera is busy. Close anything else using it and try again."
  }
  if (issue.name === "OverconstrainedError") return "Those camera settings aren't available. Try another camera."
  return issue.message || "Couldn't start the call."
}

export function friendlyCameraError(error: unknown): string {
  const issue = asError(error)
  if (issue.name === "NotFoundError" || issue.name === "OverconstrainedError") return "That camera isn't available."
  if (issue.name === "NotAllowedError") return "Camera access was denied."
  return "Couldn't switch cameras. The current camera is still selected."
}

export function friendlyScreenShareError(error: unknown): string {
  const issue = asError(error)
  if (issue.name === "NotAllowedError") return "Screen sharing wasn't allowed."
  return issue.message || "Couldn't start screen sharing."
}

export function asError(error: unknown): Error {
  return error instanceof Error ? error : new Error(String(error))
}
