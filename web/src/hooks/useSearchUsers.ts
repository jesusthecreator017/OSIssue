"use client";

import { useQuery } from "@tanstack/react-query";
import { usersApi } from "@/api/users";

export function useSearchUsers(query: string) {
	return useQuery({
		queryKey: ["users", "search", query],
		queryFn: () => usersApi.searchByName(query),
		enabled: query.length > 0,
		staleTime: 30_000,
	});
}
