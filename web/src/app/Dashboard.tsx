"use client";

import { useAuth } from "@/contexts/AuthContext";
import { useIssues } from "@/hooks/useIssues";
import { Card, CardContent } from "@/components/ui/card";

export default function Dashboard() {
	const { user } = useAuth();
	const { data: issuesData } = useIssues();

	if (!user) {
		return null;
	}

	const userIssues =
		issuesData?.filter((issue) => issue.user_id === user.id) ?? [];
	const issuesByPriority = {
		Low: userIssues.filter((i) => i.priority === "Low").length,
		Medium: userIssues.filter((i) => i.priority === "Medium").length,
		High: userIssues.filter((i) => i.priority === "High").length,
		Critical: userIssues.filter((i) => i.priority === "Critical").length,
	};

	return (
		<div>
			<h1 className="text-2xl font-bold">Welcome, {user?.name}</h1>
			<h2 className="mt-8 mb-4 text-lg font-semibold">Your Issues</h2>
			<div className="grid grid-cols-4 gap-4">
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByPriority.Low}</p>
						<p className="text-muted-foreground text-sm">Low</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByPriority.Medium}</p>
						<p className="text-muted-foreground text-sm">Medium</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByPriority.High}</p>
						<p className="text-muted-foreground text-sm">High</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByPriority.Critical}</p>
						<p className="text-muted-foreground text-sm">Critical</p>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
