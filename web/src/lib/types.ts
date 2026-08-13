export type DeviceType = "mobile" | "computer"

export interface MediaDeviceOption {
  deviceId: string
  label: string
}

export interface PeerState {
  name?: string
  device?: DeviceType
  microphoneMuted?: boolean
  noiseCancellationEnabled?: boolean
  cameraStopped?: boolean
  screenSharing?: boolean
  screenStreamID?: string
}

export interface Peer {
  id: string
  name: string
  device: DeviceType
  microphoneMuted: boolean
  noiseCancellationEnabled: boolean
  cameraStopped: boolean
  screenSharing: boolean
  connection: RTCPeerConnection
  stream: MediaStream | null
  screenStream: MediaStream | null
  cameraStreamID: string
  screenStreamID: string
  streams: Map<string, MediaStream>
  cameraSender: RTCRtpSender | null
  microphoneSender: RTCRtpSender | null
  screenSenders: RTCRtpSender[]
  locallyMuted: boolean
  renegotiate: boolean
  makingOffer: boolean
  ignoreOffer: boolean
  polite: boolean
  pendingCandidates: RTCIceCandidateInit[]
}

export interface IceServerMessage {
  urls: string | string[]
  username?: string
  credential?: string
}

declare global {
  interface Navigator {
    userAgentData?: { mobile?: boolean }
  }

  interface Window {
    webkitAudioContext?: typeof AudioContext
  }

  interface Document {
    webkitExitFullscreen?: () => void
    webkitFullscreenElement?: Element | null
  }

  interface HTMLElement {
    webkitRequestFullscreen?: () => void
  }

  interface HTMLVideoElement {
    webkitEnterFullscreen?: () => void
  }
}
