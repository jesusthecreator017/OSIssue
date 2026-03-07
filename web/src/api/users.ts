import { apiFetch } from "./client";

export interface UserSearchResult {
	id: string;
	name: string;
	email: string;
}

export const usersApi = {
	searchByName: (query: string) =>
		apiFetch<{ users: UserSearchResult[] }>(
			`/v1/users/search?q=${encodeURIComponent(query)}`,
		).then((r) => r.users),
};
