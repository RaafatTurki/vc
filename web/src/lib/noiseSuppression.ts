import { NoiseSuppressorWorklet_Name } from "@timephy/rnnoise-wasm"
import NoiseSuppressorWorklet from "@timephy/rnnoise-wasm/NoiseSuppressorWorklet?worker&url"

export class NoiseSuppression {
  private context: AudioContext | null = null
  private source: MediaStreamAudioSourceNode | null = null
  private processor: AudioWorkletNode | null = null
  private merger: ChannelMergerNode | null = null
  private destination: MediaStreamAudioDestinationNode | null = null

  async start(track: MediaStreamTrack): Promise<MediaStreamTrack> {
    await this.stop()

    const AudioContextClass = window.AudioContext || window.webkitAudioContext

    if (!AudioContextClass) throw new Error("Audio processing is not supported by this browser.")

    const context = new AudioContextClass({ sampleRate: 48000 })

    if (context.sampleRate !== 48000) {
      await context.close()
      throw new Error("This browser does not provide a 48 kHz audio context.")
    }

    await context.audioWorklet.addModule(NoiseSuppressorWorklet)

    const source = context.createMediaStreamSource(new MediaStream([track]))
    const processor = new AudioWorkletNode(context, NoiseSuppressorWorklet_Name, { channelCount: 1 })
    const merger = context.createChannelMerger(2)
    const destination = context.createMediaStreamDestination()

    source.connect(processor)
    processor.connect(merger, 0, 0)
    processor.connect(merger, 0, 1)
    merger.connect(destination)

    await context.resume()
    const processedTrack = destination.stream.getAudioTracks()[0]

    if (!processedTrack) {
      await context.close()
      throw new Error("Could not create a processed microphone track.")
    }

    this.context = context
    this.source = source
    this.processor = processor
    this.merger = merger
    this.destination = destination

    return processedTrack
  }

  async stop(): Promise<void> {

    this.processor?.disconnect()
    this.merger?.disconnect()
    this.source?.disconnect()
    this.destination?.stream.getTracks().forEach(track => track.stop())

    await this.context?.close().catch(() => {})

    this.context = null
    this.source = null
    this.processor = null
    this.merger = null
    this.destination = null
  }
}
