import { apiFetch } from "./client";
import type { Label } from "../schemas/issue";

export const labelsApi = {
	list: () => apiFetch<{ labels: Label[] }>("/v1/labels").then((r) => r.labels),

	create: (name: string, color: string) =>
		apiFetch<{ label: Label }>("/v1/labels", {
			method: "POST",
			body: JSON.stringify({ name, color }),
		}).then((r) => r.label),
};
