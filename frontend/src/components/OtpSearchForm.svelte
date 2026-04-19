<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let value = '';
	export let loading = false;

	const dispatch = createEventDispatcher<{ submit: void }>();

	function submit() {
		dispatch('submit');
	}
</script>

<form class="form" on:submit|preventDefault={submit}>
	<div class="field">
		<label for="email">Email Yahoo</label>
		<div class="input-shell">
			<span class="prefix">@</span>
			<input id="email" bind:value placeholder="user@yahoo.com" autocomplete="email" spellcheck="false" />
		</div>
	</div>
	<button type="submit" disabled={loading}>
		{#if loading}
			<span class="spinner" aria-hidden="true"></span>
			Mencari...
		{:else}
			Cari email
		{/if}
	</button>
</form>

<style>
	.form {
		display: grid;
		gap: 0.95rem;
	}
	.field {
		display: grid;
		gap: 0.45rem;
	}
	label {
		color: #cbd5e1;
		font-size: 0.92rem;
		font-weight: 600;
	}
	.input-shell {
		display: flex;
		align-items: center;
		gap: 0.65rem;
		padding: 0.95rem 1rem;
		border-radius: 18px;
		border: 1px solid rgba(148, 163, 184, 0.16);
		background: rgba(15, 23, 42, 0.82);
		transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
	}
	.input-shell:focus-within {
		border-color: rgba(96, 165, 250, 0.55);
		box-shadow: 0 0 0 4px rgba(96, 165, 250, 0.12);
		transform: translateY(-1px);
	}
	.prefix {
		color: #60a5fa;
		font-weight: 700;
		user-select: none;
	}
	input {
		width: 100%;
		border: 0;
		outline: none;
		background: transparent;
		color: #fff;
		padding: 0;
	}
	button {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.55rem;
		padding: 1rem 1rem;
		border: 0;
		border-radius: 18px;
		background: linear-gradient(135deg, #2563eb, #0ea5e9);
		color: #fff;
		font-weight: 700;
		letter-spacing: 0.01em;
		box-shadow: 0 14px 32px rgba(37, 99, 235, 0.22);
		transition: transform 0.2s ease, box-shadow 0.2s ease, opacity 0.2s ease;
	}
	button:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 18px 38px rgba(37, 99, 235, 0.28);
	}
	button:disabled { opacity: 0.72; cursor: not-allowed; }
	.spinner {
		width: 0.9rem;
		height: 0.9rem;
		border-radius: 999px;
		border: 2px solid rgba(255, 255, 255, 0.34);
		border-top-color: #fff;
		animation: spin 0.8s linear infinite;
	}
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
