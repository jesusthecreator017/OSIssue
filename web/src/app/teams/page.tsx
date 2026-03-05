"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useTeams } from "@/hooks/useTeams";
import { TeamCard } from "@/components/teams/TeamCard";
import { CreateTeamDialog } from "@/components/teams/CreateTeamDialog";

export default function TeamsPage() {
	const { user, isLoading: authLoading } = useAuth();
	const router = useRouter();
	const { data: teams, isLoading } = useTeams();

	useEffect(() => {
		if (!authLoading && !user) {
			router.push("/auth/login");
		}
	}, [authLoading, user, router]);

	if (authLoading || !user) {
		return null;
	}

	return (
		<div className="mx-auto max-w-6xl px-4 py-8">
			<div className="mb-6 flex items-center justify-between">
				<h1 className="text-2xl font-bold">Teams</h1>
				<CreateTeamDialog />
			</div>

			{isLoading ? (
				<p className="text-muted-foreground text-sm">Loading teams...</p>
			) : !teams?.length ? (
				<p className="text-muted-foreground text-sm">
					No teams yet. Create one to get started.
				</p>
			) : (
				<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
					{teams.map((team) => (
						<TeamCard key={team.id} team={team} />
					))}
				</div>
			)}
		</div>
	);
}
