export function mediaStream(node: HTMLVideoElement, stream: MediaStream | null) {
  let current = stream
  const attach = () => {
    if (node.srcObject === current) return
    node.srcObject = current
    if (current) node.play().catch(() => {})
  }
  attach()
  return { update(next: MediaStream | null) { current = next; attach() }, destroy() { node.srcObject = null } }
}
