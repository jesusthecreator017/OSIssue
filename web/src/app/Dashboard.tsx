"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useIssues } from "@/hooks/useIssues";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function Dashboard() {
	const { user } = useAuth();
	const { data: issuesData } = useIssues();

	if (!user) {
		return null;
	}

	const userIssues =
		issuesData?.filter((issue) => issue.user_id === user.id) ?? [];
	const issuesByStatus = {
		Incomplete: userIssues.filter((i) => i.status === "Incomplete").length,
		"In-Progress": userIssues.filter((i) => i.status === "In-Progress").length,
		Complete: userIssues.filter((i) => i.status === "Complete").length,
	};

	return (
		<div>
			<h1 className="text-2xl font-bold">Welcome, {user?.name}</h1>
			<h2 className="mt-8 mb-4 text-lg font-semibold">Your Issues</h2>
			<div className="grid grid-cols-3 gap-4">
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByStatus.Incomplete}</p>
						<p className="text-muted-foreground text-sm">Incomplete</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">
							{issuesByStatus["In-Progress"]}
						</p>
						<p className="text-muted-foreground text-sm">In Progress</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByStatus.Complete}</p>
						<p className="text-muted-foreground text-sm">Complete</p>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
