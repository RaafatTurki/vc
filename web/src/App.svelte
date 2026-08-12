<script lang="ts">
  import { onMount, onDestroy } from "svelte"
  import { SvelteMap } from "svelte/reactivity"
  import TopBar from "./lib/TopBar.svelte"
  import CallPage from "./pages/CallPage.svelte"
  import HomePage from "./pages/HomePage.svelte"
  import {
    cleanName,
    createRoomID,
    detectDeviceType,
    asError,
    friendlyCameraError,
    friendlyMediaError,
    friendlyScreenShareError,
    validRoom,
    videoConstraints,
  } from "./lib/call"
  import type { DeviceType, MediaDeviceOption, Peer, PeerState } from "./lib/types"

  interface SignalMessage {
    type: string
    roomId?: string
    peerId?: string
    from?: string
    peers?: string[]
    iceServers?: RTCIceServer[]
    payload?: unknown
    code?: string
    message?: string
  }

  const peers = new SvelteMap<string, Peer>()

  let roomInput = $state("")
  let signalInput = $state("")
  let participantName = $state("")
  let roomID = $state("")
  let currentPage = $state<"home" | "call">("home")
  let joining = $state(false)
  let setupError = $state("")
  let statusState = $state<"idle" | "connecting" | "connected" | "error">("idle")
  let statusText = $state("Ready")
  let localStream = $state<MediaStream | null>(null)
  let joinWithAudio = $state(true)
  let joinWithVideo = $state(true)
  let microphoneMuted = $state(false)
  let cameraStopped = $state(false)
  let cameraFacing = $state<VideoFacingModeEnum>("user")
  let canSwitchCamera = $state(false)
  let switchingCamera = $state(false)
  let audioDevices = $state<MediaDeviceOption[]>([])
  let videoDevices = $state<MediaDeviceOption[]>([])
  let selectedAudioDeviceID = $state("")
  let selectedVideoDeviceID = $state("")
  let switchingAudioDevice = $state(false)
  let switchingVideoDevice = $state(false)
  let canShareScreen = $state(false)
  let screenSharing = $state(false)
  let sharingScreen = $state(false)
  let cameraError = $state("")
  let copyLabel = $state("Copy invite link")

  let socket: WebSocket | null = null
  let selfPeerID = ""
  let iceServers: RTCIceServer[] = []
  let leaving = false
  let welcomed = false
  let signalQueue = Promise.resolve()
  let callEpoch = 0
  let deviceType = $state<DeviceType>("computer")
  let joinSound: HTMLAudioElement | null = null
  let cameraTrack: MediaStreamTrack | null = null
  let displayTrack: MediaStreamTrack | null = null
  let displayStream = $state<MediaStream | null>(null)
  let mixedAudioTrack: MediaStreamTrack | null = null
  let screenAudioContext: AudioContext | null = null
  let screenAudioNodes: AudioNode[] = []

  onMount(() => {
    const query = new URLSearchParams(window.location.search)
    const requestedRoom = query.get("room")
    signalInput = query.get("signal") || defaultSignalURL()
    participantName = localStorage.getItem("vivid-name") || ""
    deviceType = detectDeviceType()
    canShareScreen = typeof navigator.mediaDevices?.getDisplayMedia === "function"
    joinSound = new Audio("/join.opus")
    joinSound.preload = "auto"
    joinSound.volume = 0.65
    window.addEventListener("beforeunload", closeConnections)
    window.addEventListener("popstate", handleRouteChange)
    navigator.mediaDevices?.addEventListener("devicechange", refreshMediaDevices)

    const routeRoom = roomFromPath()
    if (routeRoom) {
      roomInput = routeRoom
      roomID = routeRoom
      currentPage = "call"
      if (!participantName) participantName = `Guest ${createRoomID().slice(0, 4)}`
      joinCall()
    } else {
      roomInput = validRoom(requestedRoom) ? requestedRoom : createRoomID()
      if (window.location.pathname !== "/") window.history.replaceState({}, "", "/")
    }
  })

  onDestroy(() => {
    window.removeEventListener("beforeunload", closeConnections)
    window.removeEventListener("popstate", handleRouteChange)
    navigator.mediaDevices?.removeEventListener("devicechange", refreshMediaDevices)
    closeConnections()
    joinSound = null
  })

  async function joinCall(event?: SubmitEvent): Promise<void> {
    event?.preventDefault()
    if (joining) return
    if (localStream || socket || cameraTrack || displayStream) closeConnections()
    setupError = ""
    leaving = false
    welcomed = false
    callEpoch += 1
    signalQueue = Promise.resolve()
    roomID = roomInput.trim()
    participantName = cleanName(participantName)

    if (!participantName) {
      setupError = "Enter your name before joining."
      return
    }

    if (!validRoom(roomID)) {
      setupError = "Room IDs must be exactly 6 letters or numbers."
      return
    }

    let signalURL
    try {
      signalURL = new URL(signalInput.trim())
      if (signalURL.protocol !== "ws:" && signalURL.protocol !== "wss:") {
        throw new Error("unsupported protocol")
      }
    } catch {
      setupError = "Enter a valid ws:// or wss:// signaling URL."
      return
    }

    joining = true
    setStatus("connecting", "Preparing call")

    try {
      localStream = await acquireCallMedia()
      cameraTrack = localStream.getVideoTracks()[0] || null
      microphoneMuted = !joinWithAudio
      cameraStopped = !joinWithVideo
      cameraFacing = trackFacing(localStream.getVideoTracks()[0], cameraFacing)
      await refreshMediaDevices()

      signalURL.searchParams.set("room", roomID)
      await openSignalingSocket(signalURL)
      localStorage.setItem("vivid-signal-url", signalInput.trim())
      localStorage.setItem("vivid-name", participantName)

      currentPage = "call"
      navigateToCall()
      setStatus("connecting", "Joining room")
    } catch (error) {
      closeConnections()
      joining = false
      setStatus("error", "Could not join")
      const message = friendlyMediaError(error)
      if (currentPage === "call") cameraError = message
      else setupError = message
    }
  }

  async function acquireCallMedia(): Promise<MediaStream> {
    if (!joinWithAudio && !joinWithVideo) return new MediaStream()
    try {
      return await navigator.mediaDevices.getUserMedia({
        audio: joinWithAudio,
        video: joinWithVideo ? videoConstraints(cameraFacing) : false,
      })
    } catch (error) {
      if (!["AbortError", "NotReadableError", "OverconstrainedError"].includes(asError(error).name)) {
        throw error
      }
      await new Promise(resolve => window.setTimeout(resolve, 250))
      return navigator.mediaDevices.getUserMedia({ audio: joinWithAudio, video: joinWithVideo })
    }
  }

  function openSignalingSocket(url: URL): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const candidate = new WebSocket(url)
      let settled = false

      candidate.addEventListener("open", () => {
        settled = true
        socket = candidate
        candidate.addEventListener("message", queueSignalMessage)
        candidate.addEventListener("close", onSocketClose)
        candidate.addEventListener("error", onSocketError)
        resolve()
      }, { once: true })

      candidate.addEventListener("error", () => {
        if (!settled) {
          settled = true
          reject(new Error("The signaling server could not be reached."))
        }
      }, { once: true })

      candidate.addEventListener("close", (event) => {
        if (!settled) {
          settled = true
          reject(new Error(event.reason || "The signaling server closed the connection."))
        }
      }, { once: true })
    })
  }

  function queueSignalMessage(event: MessageEvent<string>): void {
    if (event.currentTarget !== socket) return
    const epoch = callEpoch
    signalQueue = signalQueue.then(() => {
      if (epoch === callEpoch) return onSignalMessage(event, epoch)
    })
  }

  async function onSignalMessage(event: MessageEvent<string>, epoch: number): Promise<void> {
    try {
      const message = JSON.parse(event.data) as SignalMessage
      switch (message.type) {
        case "welcome":
          welcomed = true
          selfPeerID = message.peerId || ""
          iceServers = Array.isArray(message.iceServers) ? message.iceServers : []
          setStatus("connected", "Connected")
          for (const peerID of message.peers ?? []) {
            createPeer(peerID)
            sendPeerState(peerID)
            sendSignal("peer-ready", peerID, true)
          }
          break
        case "peer-joined":
          if (!message.peerId) break
          createPeer(message.peerId)
          sendPeerState(message.peerId)
          playJoinSound()
          break
        case "peer-state":
          if (message.from) receivePeerState(message.from, message.payload as PeerState)
          break
        case "peer-ready":
          if (message.from) await sendOffer(message.from)
          break
        case "peer-left":
          if (message.peerId) removePeer(message.peerId)
          break
        case "offer":
          if (message.from) await receiveOffer(message.from, message.payload as RTCSessionDescriptionInit)
          break
        case "answer":
          if (message.from) await receiveAnswer(message.from, message.payload as RTCSessionDescriptionInit)
          break
        case "ice-candidate":
          if (message.from) await receiveCandidate(message.from, message.payload as RTCIceCandidateInit)
          break
        case "error":
          setStatus("error", message.message || "Signaling error")
          break
        default:
          console.warn("Unknown signaling message", message)
      }
    } catch (error) {
      if (epoch !== callEpoch) return
      console.error("Could not process signaling message", error)
      setStatus("error", "Signaling error")
    }
  }

  function createPeer(peerID: string): Peer {
    if (!peerID) throw new Error("Peer ID is required")
    const existing = peers.get(peerID)
    if (existing) return existing
    if (!localStream) throw new Error("Local media is not ready")
    const stream = localStream

    const connection = new RTCPeerConnection({ iceServers })
    const peer: Peer = {
      id: peerID,
      name: `Guest ${peerID.slice(0, 6)}`,
      device: "computer",
      microphoneMuted: false,
      cameraStopped: false,
      screenSharing: false,
      connection,
      stream: null,
      screenStream: null,
      cameraStreamID: "",
      screenStreamID: "",
      streams: new Map(),
      cameraSender: null,
      microphoneSender: null,
      screenSenders: [],
      locallyMuted: false,
      renegotiate: false,
      makingOffer: false,
      ignoreOffer: false,
      polite: selfPeerID > peerID,
      pendingCandidates: [],
    }
    peers.set(peerID, peer)

    for (const track of stream.getTracks()) {
      const outboundTrack = track.kind === "audio" && mixedAudioTrack ? mixedAudioTrack : track
      const sender = connection.addTrack(outboundTrack, stream)
      if (track.kind === "video") peer.cameraSender = sender
      if (track.kind === "audio") peer.microphoneSender = sender
    }
    if (stream.getAudioTracks().length === 0) {
      connection.addTransceiver("audio", { direction: "recvonly" })
    }
    if (stream.getVideoTracks().length === 0) {
      connection.addTransceiver("video", { direction: "recvonly" })
    }
    if (screenSharing && displayStream && displayTrack) {
      peer.screenSenders = [connection.addTrack(displayTrack, displayStream)]
    }

    connection.addEventListener("icecandidate", ({ candidate }) => {
      if (candidate) sendSignal("ice-candidate", peerID, candidate.toJSON())
    })

    connection.addEventListener("track", (event) => {
      const currentPeer = peers.get(peerID)
      if (!currentPeer) return
      if (event.streams[0]) {
        currentPeer.streams.set(event.streams[0].id, event.streams[0])
        if (!currentPeer.cameraStreamID && (
          !currentPeer.screenSharing || event.streams[0].id !== currentPeer.screenStreamID
        )) {
          currentPeer.cameraStreamID = event.streams[0].id
        }
      } else {
        const stream = currentPeer.stream || new MediaStream()
        stream.addTrack(event.track)
        currentPeer.streams.set(stream.id, stream)
        if (!currentPeer.cameraStreamID) currentPeer.cameraStreamID = stream.id
      }
      updatePeerStreams(currentPeer)
      peers.set(peerID, { ...currentPeer })
    })

    connection.addEventListener("connectionstatechange", () => {
      if (connection.connectionState === "failed" || connection.connectionState === "closed") {
        removePeer(peerID)
      } else if (connection.connectionState === "connected") {
        setStatus("connected", `${peers.size + 1} people in call`)
      }
    })

    return peer
  }

  async function sendOffer(peerID: string): Promise<void> {
    const peer = createPeer(peerID)
    if (peer.connection.signalingState !== "stable") {
      peer.renegotiate = true
      return
    }
    peer.renegotiate = false
    peer.makingOffer = true
    try {
      const offer = await peer.connection.createOffer()
      await peer.connection.setLocalDescription(offer)
      sendSignal("offer", peerID, peer.connection.localDescription!.toJSON())
    } finally {
      peer.makingOffer = false
    }
  }

  async function receiveOffer(peerID: string, description: RTCSessionDescriptionInit): Promise<void> {
    const peer = createPeer(peerID)
    const collision = peer.makingOffer || peer.connection.signalingState !== "stable"
    peer.ignoreOffer = !peer.polite && collision
    if (peer.ignoreOffer) return
    if (collision) {
      peer.renegotiate = true
      await peer.connection.setLocalDescription({ type: "rollback" })
    }
    await peer.connection.setRemoteDescription(description)
    await flushCandidates(peer)
    const answer = await peer.connection.createAnswer()
    await peer.connection.setLocalDescription(answer)
    sendSignal("answer", peerID, peer.connection.localDescription!.toJSON())
    if (peer.renegotiate) await sendOffer(peerID)
  }

  async function receiveAnswer(peerID: string, description: RTCSessionDescriptionInit): Promise<void> {
    const peer = peers.get(peerID)
    if (!peer) throw new Error(`Answer received for unknown peer ${peerID}`)
    if (peer.connection.signalingState !== "have-local-offer") return
    await peer.connection.setRemoteDescription(description)
    await flushCandidates(peer)
    if (peer.renegotiate) await sendOffer(peerID)
  }

  async function receiveCandidate(peerID: string, candidate: RTCIceCandidateInit): Promise<void> {
    const peer = createPeer(peerID)
    try {
      if (peer.connection.remoteDescription) {
        await peer.connection.addIceCandidate(candidate)
      } else {
        peer.pendingCandidates.push(candidate)
      }
    } catch (error) {
      if (!peer.ignoreOffer) throw error
    }
  }

  async function flushCandidates(peer: Peer): Promise<void> {
    for (const candidate of peer.pendingCandidates.splice(0)) {
      await peer.connection.addIceCandidate(candidate)
    }
  }

  function sendSignal(type: string, to: string, payload: unknown): void {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type, to, payload }))
    }
  }

  function sendPeerState(peerID: string): void {
    sendSignal("peer-state", peerID, {
      name: participantName,
      device: deviceType,
      microphoneMuted,
      cameraStopped,
      screenSharing,
      screenStreamID: displayStream?.id || "",
    })
  }

  function playJoinSound(): void {
    if (!joinSound) return
    joinSound.currentTime = 0
    joinSound.play().catch(() => {})
  }

  function broadcastPeerState(): void {
    for (const peerID of peers.keys()) sendPeerState(peerID)
  }

  function receivePeerState(peerID: string, state: PeerState): void {
    const peer = createPeer(peerID)
    peer.name = cleanName(state?.name) || peer.name
    peer.device = state?.device === "mobile" ? "mobile" : "computer"
    peer.microphoneMuted = state?.microphoneMuted === true
    peer.cameraStopped = state?.cameraStopped === true
    peer.screenSharing = state?.screenSharing === true
    peer.screenStreamID = typeof state?.screenStreamID === "string" ? state.screenStreamID.slice(0, 128) : ""
    updatePeerStreams(peer)
    peers.set(peerID, { ...peer })
  }

  function updatePeerStreams(peer: Peer): void {
    const streams = [...peer.streams.values()]
    const claimedScreen = streams.find(({ id }) => id === peer.screenStreamID) || null
    peer.stream = streams.find(({ id }) => id === peer.cameraStreamID)
      || streams.find(stream => stream !== claimedScreen && stream.getAudioTracks().length > 0)
      || streams.find(stream => stream !== claimedScreen)
      || null
    if (peer.stream && !peer.cameraStreamID) peer.cameraStreamID = peer.stream.id
    peer.screenStream = peer.screenSharing
      ? claimedScreen || streams.find(stream => stream !== peer.stream && stream.getVideoTracks().length > 0) || null
      : null
  }

  function togglePeerPlayback(peerID: string): void {
    const peer = peers.get(peerID)
    if (!peer) return
    peer.locallyMuted = !peer.locallyMuted
    peers.set(peerID, { ...peer })
  }

  function removePeer(peerID: string): void {
    const peer = peers.get(peerID)
    if (!peer) return
    peers.delete(peerID)
    peer.connection.close()
    setStatus("connected", peers.size ? `${peers.size + 1} people in call` : "Connected")
  }

  function leaveCall(): void {
    goHome(true)
  }

  function goHome(pushHistory: boolean): void {
    leaving = true
    closeConnections()
    welcomed = false
    selfPeerID = ""
    joining = false
    currentPage = "home"
    roomInput = roomID
    microphoneMuted = false
    cameraStopped = false
    screenSharing = false
    sharingScreen = false
    cameraError = ""
    setStatus("idle", "Ready to rejoin")
    if (pushHistory) window.history.pushState({}, "", "/")
  }

  function closeConnections(): void {
    callEpoch += 1
    for (const peerID of [...peers.keys()]) removePeer(peerID)

    const currentSocket = socket
    socket = null
    if (currentSocket && currentSocket.readyState < WebSocket.CLOSING) {
      currentSocket.close(1000, "left call")
    }

    if (displayTrack) displayTrack.onended = null
    if (displayStream) {
      for (const track of displayStream.getTracks()) track.stop()
    }
    if (localStream) {
      for (const track of localStream.getTracks()) track.stop()
      localStream = null
    }
    cameraTrack?.stop()
    displayTrack?.stop()
    mixedAudioTrack?.stop()
    screenAudioContext?.close().catch(() => {})
    cameraTrack = null
    displayTrack = null
    displayStream = null
    mixedAudioTrack = null
    screenAudioContext = null
    screenAudioNodes = []
    screenSharing = false
    sharingScreen = false
  }

  function onSocketClose(event: CloseEvent): void {
    if (event.currentTarget !== socket) return
    if (leaving) return
    if (!welcomed) {
      closeConnections()
      joining = false
      const message = event.reason || "The signaling server closed the connection."
      if (currentPage === "call") cameraError = message
      else setupError = message
    }
    setStatus("error", event.reason || "Disconnected")
  }

  function onSocketError(event: Event): void {
    if (event.currentTarget !== socket) return
    if (!leaving) setStatus("error", "Connection problem")
  }

  async function toggleMicrophone(): Promise<void> {
    const track = localStream?.getAudioTracks()[0]
    if (!track) {
      await enableMicrophone()
      return
    }
    track.enabled = !track.enabled
    microphoneMuted = !track.enabled
    broadcastPeerState()
  }

  async function toggleCamera(): Promise<void> {
    const track = cameraTrack
    if (!track) {
      await enableCamera()
      return
    }
    track.enabled = !track.enabled
    cameraStopped = !track.enabled
    broadcastPeerState()
  }

  async function enableMicrophone(): Promise<void> {
    cameraError = ""
    try {
      const audio = selectedAudioDeviceID
        ? { deviceId: { exact: selectedAudioDeviceID } }
        : true
      const stream = await navigator.mediaDevices.getUserMedia({ audio, video: false })
      const track = stream.getAudioTracks()[0]
      if (!track) throw new Error("No microphone was available.")
      if (!localStream) throw new Error("The call is no longer active.")
      track.enabled = true
      localStream = new MediaStream([track, ...localStream.getVideoTracks()])
      let addedSender = false
      for (const peer of peers.values()) {
        if (peer.microphoneSender) {
          await peer.microphoneSender.replaceTrack(track)
        } else {
          peer.microphoneSender = peer.connection.addTrack(track, localStream)
          addedSender = true
        }
      }
      if (screenSharing && displayStream) await startScreenAudioMix(displayStream)
      microphoneMuted = false
      broadcastPeerState()
      if (addedSender) await renegotiatePeers()
      await refreshMediaDevices()
    } catch (error) {
      cameraError = asError(error).message || "Could not start the microphone."
    }
  }

  async function enableCamera(): Promise<void> {
    cameraError = ""
    try {
      const video = selectedVideoDeviceID
        ? { width: { ideal: 1280 }, height: { ideal: 720 }, deviceId: { exact: selectedVideoDeviceID } }
        : videoConstraints(cameraFacing)
      const stream = await navigator.mediaDevices.getUserMedia({ audio: false, video })
      const track = stream.getVideoTracks()[0]
      if (!track) throw new Error("No camera was available.")
      if (!localStream) throw new Error("The call is no longer active.")
      track.enabled = true
      localStream = new MediaStream([...localStream.getAudioTracks(), track])
      cameraTrack = track
      let addedSender = false
      for (const peer of peers.values()) {
        if (peer.cameraSender) {
          await peer.cameraSender.replaceTrack(track)
        } else {
          peer.cameraSender = peer.connection.addTrack(track, localStream)
          addedSender = true
        }
      }
      cameraFacing = trackFacing(track, cameraFacing)
      cameraStopped = false
      broadcastPeerState()
      if (addedSender) await renegotiatePeers()
      await refreshMediaDevices()
    } catch (error) {
      cameraError = friendlyCameraError(error)
    }
  }

  async function switchCamera(): Promise<void> {
    const oldTrack = cameraTrack
    if (!oldTrack || switchingCamera) return

    switchingCamera = true
    cameraError = ""
    const previousFacing = cameraFacing
    const nextFacing = previousFacing === "environment" ? "user" : "environment"
    const wasEnabled = oldTrack.enabled
    let oldTrackStopped = false

    try {
      let cameraStream
      try {
        cameraStream = await navigator.mediaDevices.getUserMedia({
          audio: false,
          video: videoConstraints(nextFacing, true),
        })
      } catch (error) {
        const issue = asError(error)
        if (issue.name !== "NotReadableError" && issue.name !== "AbortError") throw error
        oldTrack.stop()
        oldTrackStopped = true
        cameraStream = await navigator.mediaDevices.getUserMedia({
          audio: false,
          video: videoConstraints(nextFacing, true),
        })
      }

      const newTrack = cameraStream.getVideoTracks()[0]
      if (!newTrack) throw new Error("The selected camera did not provide video.")
      newTrack.enabled = wasEnabled
      await replaceVideoTrack(newTrack)
      cameraTrack = newTrack
      if (!oldTrackStopped) oldTrack.stop()
      cameraFacing = trackFacing(newTrack, nextFacing)
      cameraStopped = !newTrack.enabled
      await refreshMediaDevices()
    } catch (error) {
      if (oldTrackStopped) await restoreCamera(previousFacing, wasEnabled)
      cameraError = friendlyCameraError(error)
    } finally {
      switchingCamera = false
    }
  }

  async function replaceVideoTrack(newTrack: MediaStreamTrack): Promise<void> {
    if (!localStream) throw new Error("The call is no longer active.")
    const replacements = []
    let addedSender = false
    localStream = new MediaStream([
      ...localStream.getAudioTracks(),
      newTrack,
    ])
    for (const peer of peers.values()) {
      if (peer.cameraSender) {
        replacements.push(peer.cameraSender.replaceTrack(newTrack))
      } else {
        peer.cameraSender = peer.connection.addTrack(newTrack, localStream)
        addedSender = true
      }
    }
    await Promise.all(replacements)
    if (addedSender) await renegotiatePeers()
  }

  async function restoreCamera(facingMode: VideoFacingModeEnum, enabled: boolean): Promise<void> {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: false,
        video: videoConstraints(facingMode),
      })
      const track = stream.getVideoTracks()[0]
      if (!track) return
      track.enabled = enabled
      await replaceVideoTrack(track)
      cameraTrack = track
      cameraFacing = trackFacing(track, facingMode)
    } catch (error) {
      console.error("Could not restore the previous camera", error)
    }
  }

  async function toggleScreenShare(): Promise<void> {
    if (screenSharing) {
      await stopScreenShare()
    } else {
      await startScreenShare()
    }
  }

  async function startScreenShare(): Promise<void> {
    if (!canShareScreen || sharingScreen || !localStream) return
    sharingScreen = true
    cameraError = ""
    const epoch = callEpoch
    let track = null
    let sharedStream = null

    try {
      const displayOptions: DisplayMediaStreamOptions & {
        systemAudio: "include"
        windowAudio: "system"
      } = {
        video: true,
        audio: {
          echoCancellation: false,
          noiseSuppression: false,
          autoGainControl: false,
        },
        systemAudio: "include",
        windowAudio: "system",
      }
      sharedStream = await navigator.mediaDevices.getDisplayMedia(displayOptions)
      track = sharedStream.getVideoTracks()[0]
      if (!track) throw new Error("Screen sharing did not provide video.")
      if (epoch !== callEpoch || !localStream) {
        for (const sharedTrack of sharedStream.getTracks()) sharedTrack.stop()
        return
      }

      if (track.readyState === "ended") throw new Error("Screen sharing ended before it started.")
      displayStream = sharedStream
      displayTrack = track
      screenSharing = true
      track.onended = () => stopScreenShare()
      const hasScreenAudio = await startScreenAudioMix(sharedStream)
      for (const peer of peers.values()) {
        peer.screenSenders = [peer.connection.addTrack(track, sharedStream)]
      }
      broadcastPeerState()
      await renegotiatePeers()
      if (!hasScreenAudio) {
        cameraError = "This browser did not provide screen audio. Enable ‘Share audio’ in the sharing dialog if available."
      }
    } catch (error) {
      for (const peer of peers.values()) {
        for (const sender of peer.screenSenders) peer.connection.removeTrack(sender)
        peer.screenSenders = []
      }
      for (const sharedTrack of sharedStream?.getTracks() || []) sharedTrack.stop()
      await stopScreenAudioMix()
      displayTrack = null
      displayStream = null
      screenSharing = false
      broadcastPeerState()
      if (epoch === callEpoch) cameraError = friendlyScreenShareError(error)
    } finally {
      if (epoch === callEpoch) sharingScreen = false
    }
  }

  async function stopScreenShare(): Promise<void> {
    if (!screenSharing || sharingScreen) return
    sharingScreen = true
    cameraError = ""
    const track = displayTrack
    const stream = displayStream
    displayTrack = null
    displayStream = null
    screenSharing = false
    if (track) track.onended = null

    try {
      await stopScreenAudioMix()
      for (const peer of peers.values()) {
        for (const sender of peer.screenSenders) peer.connection.removeTrack(sender)
        peer.screenSenders = []
      }
      broadcastPeerState()
      await renegotiatePeers()
    } catch (error) {
      cameraError = asError(error).message || "Could not stop screen sharing cleanly."
    } finally {
      for (const sharedTrack of stream?.getTracks() || []) sharedTrack.stop()
      sharingScreen = false
    }
  }

  async function renegotiatePeers(): Promise<void> {
    for (const peerID of peers.keys()) await sendOffer(peerID)
  }

  async function startScreenAudioMix(sharedStream: MediaStream): Promise<boolean> {
    const microphoneTrack = localStream?.getAudioTracks()[0]
    const screenAudioTrack = sharedStream.getAudioTracks()[0]
    if (!screenAudioTrack) return false

    const AudioContextClass = window.AudioContext || window.webkitAudioContext
    if (!AudioContextClass) return false

    const context = new AudioContextClass()
    const destination = context.createMediaStreamDestination()
    const screenSource = context.createMediaStreamSource(new MediaStream([screenAudioTrack]))
    screenSource.connect(destination)
    const microphoneSource = microphoneTrack
      ? context.createMediaStreamSource(new MediaStream([microphoneTrack]))
      : null
    microphoneSource?.connect(destination)
    await context.resume()

    const track = destination.stream.getAudioTracks()[0]
    if (!track) {
      await context.close()
      return false
    }

    screenAudioContext = context
    screenAudioNodes = [microphoneSource, screenSource, destination].filter(
      (node): node is MediaStreamAudioSourceNode | MediaStreamAudioDestinationNode => node !== null,
    )
    mixedAudioTrack = track
    await replaceAudioTrack(track)
    return true
  }

  async function stopScreenAudioMix(replaceSender = true): Promise<void> {
    const microphoneTrack = localStream?.getAudioTracks()[0] || null
    if (replaceSender && mixedAudioTrack) {
      await Promise.all([...peers.values()].map(peer => (
        peer.microphoneSender?.replaceTrack(microphoneTrack)
      )))
    }
    mixedAudioTrack?.stop()
    mixedAudioTrack = null
    screenAudioNodes = []
    const context = screenAudioContext
    screenAudioContext = null
    await context?.close().catch(() => {})
  }

  async function changeAudioDevice(deviceID: string): Promise<void> {
    if (!deviceID || switchingAudioDevice || !localStream) return
    switchingAudioDevice = true
    cameraError = ""
    let newTrack = null
    let oldTrack = null

    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: { deviceId: { exact: deviceID } },
        video: false,
      })
      newTrack = stream.getAudioTracks()[0]
      if (!newTrack) throw new Error("The selected microphone did not provide audio.")

      oldTrack = localStream.getAudioTracks()[0]
      newTrack.enabled = oldTrack?.enabled ?? true
      if (mixedAudioTrack) await stopScreenAudioMix(false)
      localStream = new MediaStream([newTrack, ...localStream.getVideoTracks()])

      let mixed = false
      if (screenSharing && displayStream) mixed = await startScreenAudioMix(displayStream)
      if (!mixed) await replaceAudioTrack(newTrack)
      oldTrack?.stop()
      selectedAudioDeviceID = newTrack.getSettings().deviceId || deviceID
      await refreshMediaDevices()
    } catch (error) {
      if (mixedAudioTrack) await stopScreenAudioMix(false)
      if (newTrack && localStream?.getTracks().includes(newTrack)) {
        await replaceAudioTrack(newTrack).catch(() => {})
        oldTrack?.stop()
      }
      if (newTrack && !localStream?.getTracks().includes(newTrack)) newTrack.stop()
      cameraError = asError(error).message || "Could not switch microphones."
    } finally {
      switchingAudioDevice = false
    }
  }

  async function changeVideoDevice(deviceID: string): Promise<void> {
    if (!deviceID || switchingVideoDevice || !localStream) return
    switchingVideoDevice = true
    cameraError = ""
    let newTrack = null

    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: false,
        video: {
          width: { ideal: 1280 },
          height: { ideal: 720 },
          deviceId: { exact: deviceID },
        },
      })
      newTrack = stream.getVideoTracks()[0]
      if (!newTrack) throw new Error("The selected camera did not provide video.")

      const oldTrack = cameraTrack
      newTrack.enabled = oldTrack?.enabled ?? true
      await replaceVideoTrack(newTrack)
      cameraTrack = newTrack
      oldTrack?.stop()
      cameraFacing = trackFacing(newTrack, cameraFacing)
      selectedVideoDeviceID = newTrack.getSettings().deviceId || deviceID
      await refreshMediaDevices()
    } catch (error) {
      if (newTrack && newTrack !== cameraTrack) newTrack.stop()
      cameraError = friendlyCameraError(error)
    } finally {
      switchingVideoDevice = false
    }
  }

  async function replaceAudioTrack(newTrack: MediaStreamTrack | null): Promise<void> {
    const replacements = []
    let addedSender = false
    for (const peer of peers.values()) {
      if (peer.microphoneSender) {
        replacements.push(peer.microphoneSender.replaceTrack(newTrack))
      } else if (newTrack && localStream) {
        peer.microphoneSender = peer.connection.addTrack(newTrack, localStream)
        addedSender = true
      }
    }
    await Promise.all(replacements)
    if (addedSender) await renegotiatePeers()
  }

  async function refreshMediaDevices(): Promise<void> {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices()
      audioDevices = devices
        .filter(({ kind }) => kind === "audioinput")
        .map(({ deviceId, label }, index) => ({ deviceId, label: label || `Microphone ${index + 1}` }))
      videoDevices = devices
        .filter(({ kind }) => kind === "videoinput")
        .map(({ deviceId, label }, index) => ({ deviceId, label: label || `Camera ${index + 1}` }))
      canSwitchCamera = videoDevices.length > 1
      selectedAudioDeviceID = localStream?.getAudioTracks()[0]?.getSettings().deviceId || selectedAudioDeviceID
      selectedVideoDeviceID = cameraTrack?.getSettings().deviceId || selectedVideoDeviceID
    } catch {
      audioDevices = []
      videoDevices = []
      canSwitchCamera = false
    }
  }

  function trackFacing(track: MediaStreamTrack | undefined, fallback: VideoFacingModeEnum): VideoFacingModeEnum {
    const facingMode = track?.getSettings().facingMode
    if (facingMode === "user" || facingMode === "environment" || facingMode === "left" || facingMode === "right") {
      return facingMode
    }
    return fallback
  }

  async function copyInviteLink(): Promise<void> {
    try {
      await navigator.clipboard.writeText(inviteURL())
      copyLabel = "Copied"
    } catch {
      window.prompt("Copy this invite link:", inviteURL())
    }
    window.setTimeout(() => copyLabel = "Copy invite link", 1600)
  }

  function inviteURL(): string {
    const url = new URL(window.location.href)
    url.search = ""
    url.hash = ""
    url.pathname = `/${encodeURIComponent(roomID)}`
    return url.toString()
  }

  function navigateToCall(): void {
    const nextPath = `/${encodeURIComponent(roomID)}`
    if (window.location.pathname !== nextPath || window.location.search || window.location.hash) {
      window.history.pushState({}, "", nextPath)
    }
  }

  function roomFromPath(): string {
    const parts = window.location.pathname.split("/").filter(Boolean)
    if (parts.length !== 1) return ""
    try {
      const value = decodeURIComponent(parts[0])
      return validRoom(value) ? value : ""
    } catch {
      return ""
    }
  }

  function handleRouteChange(): void {
    const nextRoom = roomFromPath()
    if (!nextRoom) {
      goHome(false)
      return
    }
    if (currentPage === "call" && roomID === nextRoom) return

    leaving = true
    closeConnections()
    leaving = false
    welcomed = false
    joining = false
    cameraError = ""
    roomInput = nextRoom
    roomID = nextRoom
    currentPage = "call"
    if (!participantName) participantName = `Guest ${createRoomID().slice(0, 4)}`
    joinCall()
  }

  function setStatus(state: "idle" | "connecting" | "connected" | "error", text: string): void {
    statusState = state
    statusText = text
  }

  function defaultSignalURL(): string {
    const saved = localStorage.getItem("vivid-signal-url")
    if (saved) return saved
    if (window.location.hostname === "raafat.io" || window.location.hostname.endsWith(".raafat.io")) {
      return "wss://signal.raafat.io/v1/ws"
    }
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    return `${protocol}//${window.location.hostname || "localhost"}:8080/v1/ws`
  }

</script>

<svelte:head><title>{currentPage === "call" ? `${roomID} · Vivid` : "Vivid"}</title></svelte:head>

<main class="shell">
  <TopBar state={statusState} text={statusText} />

  {#if currentPage === "home"}
    <HomePage
      bind:roomInput
      bind:signalInput
      bind:participantName
      bind:joinWithAudio
      bind:joinWithVideo
      {joining}
      {setupError}
      onJoin={joinCall}
    />
  {:else}
    <CallPage
      {roomID}
      {copyLabel}
      {localStream}
      {participantName}
      {deviceType}
      {microphoneMuted}
      {cameraStopped}
      {cameraFacing}
      {screenSharing}
      {displayStream}
      {peers}
      {audioDevices}
      {videoDevices}
      bind:selectedAudioDeviceID
      bind:selectedVideoDeviceID
      {switchingAudioDevice}
      {switchingVideoDevice}
      {canShareScreen}
      {sharingScreen}
      {canSwitchCamera}
      {switchingCamera}
      {cameraError}
      onCopy={copyInviteLink}
      onTogglePeerPlayback={togglePeerPlayback}
      onAudioDeviceChange={changeAudioDevice}
      onVideoDeviceChange={changeVideoDevice}
      onToggleMicrophone={toggleMicrophone}
      onToggleCamera={toggleCamera}
      onToggleScreenShare={toggleScreenShare}
      onSwitchCamera={switchCamera}
      onLeave={leaveCall}
    />
  {/if}
</main>

<style>
  .shell {
    width: min(var(--content-width), calc(100% - 2.5rem));
    min-height: 100vh;
    margin: 0 auto;
    padding: 2.25rem 0 6rem;
  }

  @media (max-width: 47.5em) {
    .shell {
      width: min(var(--content-width), calc(100% - 2rem));
      padding-top: 1.75rem;
      padding-bottom: 4.5rem;
    }
  }
</style>
