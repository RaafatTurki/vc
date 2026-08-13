<script lang="ts">
  import { MessageCircle, Send, X } from "@lucide/svelte"
  import { tick } from "svelte"

  export interface ChatMessage {
    id: string
    senderID: string
    senderName: string
    text: string
    own: boolean
  }

  interface Props {
    messages: ChatMessage[]
    open: boolean
    unread: number
    onSend: (text: string) => void
    onToggle: () => void
  }

  let { messages, open, unread, onSend, onToggle }: Props = $props()
  let text = $state("")
  let messagesElement = $state<HTMLDivElement>()

  $effect(() => {
    messages.length
    if (!open) return
    tick().then(() => {
      if (messagesElement) messagesElement.scrollTop = messagesElement.scrollHeight
    })
  })

  function send(event: SubmitEvent) {
    event.preventDefault()
    const value = text.trim()
    if (!value) return
    onSend(value)
    text = ""
  }
</script>

<button class="chat-toggle" type="button" aria-label={open ? "Close chat" : "Open chat"} aria-expanded={open} onclick={onToggle}>
  {#if open}<X aria-hidden="true" />{:else}<MessageCircle aria-hidden="true" />{/if}
  {#if unread && !open}<span class="unread">{unread > 99 ? "99+" : unread}</span>{/if}
</button>

{#if open}
  <aside class="chat-panel" aria-label="Room chat">
    <header><strong>Chat</strong><span>{messages.length} messages</span></header>
    <div bind:this={messagesElement} class="messages" aria-live="polite">
      {#if messages.length === 0}<p class="empty">No messages yet.</p>{/if}
      {#each messages as message (message.id)}
        <article class:own={message.own}>
          <strong>{message.own ? "You" : message.senderName}</strong>
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
  .chat-toggle { position: fixed; right: 1.25rem; bottom: 1.25rem; z-index: 30; display: grid; width: var(--control-height); height: var(--control-height); place-items: center; border: 1px solid var(--accent-border); border-radius: 2px; color: var(--ink); background: var(--surface); }
  .chat-toggle :global(svg) { width: 1.1rem; height: 1.1rem; }
  .unread { position: absolute; top: -0.45rem; right: -0.45rem; min-width: 1.15rem; padding: 0.1rem; border-radius: 999px; color: var(--bg); background: var(--accent); font-size: 0.65rem; text-align: center; }
  .chat-panel { position: fixed; right: 1.25rem; bottom: 4.5rem; z-index: 30; display: grid; grid-template-rows: auto 1fr auto; width: min(22rem, calc(100vw - 2rem)); height: min(32rem, calc(100vh - 7rem)); border: 1px solid var(--accent-border); border-radius: 2px; background: var(--surface); box-shadow: 0 1rem 3rem var(--accent-16); }
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
  @media (max-width: 38.75em) { .chat-toggle { right: 0.75rem; bottom: 0.75rem; } .chat-panel { right: 0.75rem; bottom: 4rem; width: calc(100vw - 1.5rem); height: min(30rem, calc(100vh - 5rem)); } }
</style>
