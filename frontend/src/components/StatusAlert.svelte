<script>
	export let status = '';
	export let message = '';

	const tones = {
		found: 'success',
		not_found: 'neutral',
		expired: 'warning',
		invalid_email: 'danger',
		error: 'danger'
	};
	const tone = tones[status] ?? 'neutral';
</script>

<div class={`alert ${tone}`}>
	<div class="dot"></div>
	<div>
		<strong>{status}</strong>
		<span>{message}</span>
	</div>
</div>

<style>
	.alert {
		padding: 0.85rem 0.95rem;
		border-radius: 16px;
		display: flex;
		align-items: start;
		gap: 0.75rem;
		margin-top: 1rem;
		font-size: 0.92rem;
		backdrop-filter: blur(10px);
	}
	.success { background: rgba(52, 211, 153, 0.1); border: 1px solid rgba(52, 211, 153, 0.22); }
	.warning { background: rgba(251, 191, 36, 0.1); border: 1px solid rgba(251, 191, 36, 0.22); }
	.danger { background: rgba(248, 113, 113, 0.1); border: 1px solid rgba(248, 113, 113, 0.22); }
	.neutral { background: rgba(148, 163, 184, 0.08); border: 1px solid rgba(148, 163, 184, 0.16); }
	.alert strong {
		display: block;
		text-transform: capitalize;
	}
	.alert span {
		display: block;
		color: #cbd5e1;
		margin-top: 0.1rem;
	}
	.dot {
		width: 0.7rem;
		height: 0.7rem;
		border-radius: 999px;
		margin-top: 0.25rem;
		background: currentColor;
		flex: 0 0 auto;
		opacity: 0.9;
	}
</style>
