import { expect, test } from "bun:test"
import { cleanName, createRoomID, validRoom, videoConstraints } from "./call"

test("room IDs are six URL-safe alphanumeric characters", () => {
  const room = createRoomID()
  expect(validRoom(room)).toBe(true)
  expect(validRoom("short")).toBe(false)
  expect(validRoom("abcdef!")).toBe(false)
})

test("names remove control characters and respect the display limit", () => {
  expect(cleanName(" \u0000 Ada\n ")).toBe("Ada")
  expect(cleanName("x".repeat(40))).toHaveLength(32)
})

test("camera constraints preserve the requested facing policy", () => {
  expect(videoConstraints("user").facingMode).toEqual({ ideal: "user" })
  expect(videoConstraints("environment", true).facingMode).toEqual({ exact: "environment" })
})
