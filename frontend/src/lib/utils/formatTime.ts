export function formatTime(value?: string): string {
	if (!value) return '-';
	return new Date(value).toLocaleString();
}
