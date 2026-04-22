<script>
	import { formatTime } from '$lib/utils/formatTime.js';

	export let result = null;

	const statusLabel = {
		found: 'Found',
		not_found: 'Not found',
		expired: 'Expired',
		invalid_email: 'Invalid email',
		error: 'Error'
	};

	const statusTone = {
		found: 'success',
		not_found: 'neutral',
		expired: 'warning',
		invalid_email: 'danger',
		error: 'danger'
	};
</script>

{#if result}
<section class="card">
	<div class="header">
		<div>
			<p class="eyebrow">Lookup result</p>
			<h3>Full text</h3>
			<p class="summary">{result.message}</p>
		</div>
		<span class={`badge ${statusTone[result.status] ?? 'neutral'}`}>
			{statusLabel[result.status] ?? result.status}
		</span>
	</div>

	<div class="meta-grid">
		<div class="meta">
			<span class="label">Email</span>
			<strong>{result.email}</strong>
		</div>
		<div class="meta">
			<span class="label">Diterima</span>
			<strong>{formatTime(result.receivedAt)}</strong>
		</div>
	</div>

	{#if result.text}
		<div class="body-block">
			<span class="label">Full text</span>
			<pre>{result.text}</pre>
		</div>
	{/if}
</section>
{/if}

<style>
	.card {
		margin-top: 1rem;
		padding: 1.1rem;
		border-radius: 28px;
		background: linear-gradient(180deg, rgba(8, 15, 27, 0.92), rgba(8, 15, 27, 0.76));
		border: 1px solid rgba(148, 163, 184, 0.14);
		box-shadow: 0 28px 70px rgba(2, 6, 23, 0.28);
	}
	.header {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		align-items: start;
		padding-bottom: 1rem;
		border-bottom: 1px solid rgba(148, 163, 184, 0.12);
	}
	.eyebrow {
		margin: 0 0 0.3rem;
		text-transform: uppercase;
		letter-spacing: 0.16em;
		font-size: 0.72rem;
		color: #93c5fd;
		font-weight: 700;
	}
	h3 {
		margin: 0;
		font-size: 1.1rem;
		letter-spacing: -0.02em;
	}
	.summary {
		margin: 0.45rem 0 0;
		color: #cbd5e1;
	}
	.badge {
		padding: 0.45rem 0.75rem;
		border-radius: 999px;
		font-size: 0.78rem;
		font-weight: 700;
		border: 1px solid transparent;
		white-space: nowrap;
	}
	.badge.success { background: rgba(52, 211, 153, 0.14); border-color: rgba(52, 211, 153, 0.26); color: #86efac; }
	.badge.warning { background: rgba(251, 191, 36, 0.14); border-color: rgba(251, 191, 36, 0.26); color: #fde68a; }
	.badge.danger { background: rgba(248, 113, 113, 0.14); border-color: rgba(248, 113, 113, 0.26); color: #fca5a5; }
	.badge.neutral { background: rgba(148, 163, 184, 0.12); border-color: rgba(148, 163, 184, 0.2); color: #e2e8f0; }
	.meta-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 0.85rem;
		margin-top: 1rem;
	}
	.meta {
		padding: 0.95rem 1rem;
		border-radius: 18px;
		background: rgba(15, 23, 42, 0.55);
		border: 1px solid rgba(148, 163, 184, 0.12);
	}
	.meta.accent {
		background: linear-gradient(135deg, rgba(37, 99, 235, 0.18), rgba(34, 197, 94, 0.12));
	}
	.meta .label {
		display: block;
		margin-bottom: 0.3rem;
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #94a3b8;
	}
	.meta strong {
		font-weight: 600;
		line-height: 1.45;
		word-break: break-word;
	}
	.body-block {
		margin-top: 1rem;
		padding: 0.95rem 1rem;
		border-radius: 18px;
		background: rgba(15, 23, 42, 0.55);
		border: 1px solid rgba(148, 163, 184, 0.12);
	}
	.body-block .label {
		display: block;
		margin-bottom: 0.55rem;
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #94a3b8;
	}
	.body-block pre {
		margin: 0;
		white-space: pre-wrap;
		word-break: break-word;
		font-family: inherit;
		line-height: 1.6;
		color: #e2e8f0;
	}
	@media (max-width: 640px) {
		.meta-grid { grid-template-columns: 1fr; }
		.header { flex-direction: column; }
	}
</style>
