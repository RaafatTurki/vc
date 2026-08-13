<script lang="ts">
  import { Send } from "@lucide/svelte"
  import { tick } from "svelte"

  export interface ChatMessage {
    id: string
    senderID: string
    senderName: string
    text: string
    timestamp: number
    own: boolean
  }

  interface Props {
    messages: ChatMessage[]
    open: boolean
    unread: number
    onSend: (text: string) => void
  }

  let { messages, open, unread, onSend }: Props = $props()
  let text = $state("")
  let messagesElement = $state<HTMLDivElement>()

  $effect(() => {
    messages.length
    if (!open) return
    tick().then(() => {
      if (messagesElement) messagesElement.scrollTop = messagesElement.scrollHeight
    })
  })

  function formatTime(timestamp: number): string {
    return new Intl.DateTimeFormat(undefined, { hour: "2-digit", minute: "2-digit", hour12: false }).format(timestamp)
  }

  function send(event: SubmitEvent) {
    event.preventDefault()
    const value = text.trim()
    if (!value) return
    onSend(value)
    text = ""
  }
</script>

{#if open}
  <aside class="chat-panel" aria-label="Room chat">
    <header><strong>Chat</strong><span>{messages.length} messages</span></header>
    <div bind:this={messagesElement} class="messages" aria-live="polite">
      {#if messages.length === 0}<p class="empty">No messages yet.</p>{/if}
      {#each messages as message (message.id)}
        <article class:own={message.own}>
          <strong>{message.own ? "You" : message.senderName}</strong>
          <time datetime={new Date(message.timestamp).toISOString()}>{formatTime(message.timestamp)}</time>
          <p>{message.text}</p>
        </article>
      {/each}
    </div>
    <form onsubmit={send}>
      <input bind:value={text} maxlength="4000" placeholder="Write a message" aria-label="Message" />
      <button type="submit" aria-label="Send message"><Send aria-hidden="true" /></button>
    </form>
  </aside>
{/if}

<style>
  .chat-panel { display: grid; grid-template-rows: auto minmax(0, 1fr) auto; width: 100%; height: 40rem; max-height: 75vh; border: 1px solid var(--accent-border); border-radius: 2px; background: var(--surface); }
  header { display: flex; justify-content: space-between; padding: var(--space-3); border-bottom: 1px solid var(--line-soft); }
  header span { color: var(--muted); font-size: 0.72rem; }
  .messages { overflow-y: auto; padding: var(--space-3); }
  article { width: fit-content; max-width: 90%; margin-bottom: var(--space-3); padding: 0.45rem 0.6rem; border: 1px solid var(--line-soft); background: var(--surface-faint); }
  article.own { border-color: var(--accent-border); background: var(--accent-subtle); }
  article strong { color: var(--accent); font-size: 0.7rem; }
  article p { margin: 0.2rem 0 0; overflow-wrap: anywhere; white-space: pre-wrap; }
  .empty { color: var(--muted); font-size: 0.82rem; text-align: center; }
  form { display: flex; gap: var(--space-2); padding: var(--space-3); border-top: 1px solid var(--line-soft); }
  input { min-width: 0; flex: 1; height: var(--control-height); padding: 0 0.6rem; border: 1px solid var(--line-soft); color: var(--ink); background: var(--surface-faint); }
  form button { display: grid; width: var(--control-height); height: var(--control-height); place-items: center; border: 1px solid var(--accent-border); color: var(--ink); background: var(--accent-subtle); }
  form button :global(svg) { width: 1rem; height: 1rem; }
  time { margin-left: var(--space-2); color: var(--muted); font-size: 0.65rem; }
  @media (max-width: 47.5em) { .chat-panel { position: static; width: 100%; height: min(28rem, 55vh); } }
</style>
