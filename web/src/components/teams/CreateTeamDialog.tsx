"use client";

import { useState } from "react";
import { useForm } from "@tanstack/react-form";
import { Plus } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useCreateTeam } from "@/hooks/useTeams";

export function CreateTeamDialog() {
	const [open, setOpen] = useState(false);
	const createTeam = useCreateTeam();

	const form = useForm({
		defaultValues: {
			name: "",
			description: "",
			max_members: 50,
		},
		onSubmit: async ({ value }) => {
			try {
				await createTeam.mutateAsync(value);
				toast.success("Team created");
				setOpen(false);
				form.reset();
			} catch (err) {
				toast.error(err instanceof Error ? err.message : "Failed to create team");
			}
		},
	});

	return (
		<Dialog open={open} onOpenChange={setOpen}>
			<DialogTrigger asChild>
				<Button>
					<Plus className="mr-2 size-4" />
					Create Team
				</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Create Team</DialogTitle>
				</DialogHeader>
				<form
					onSubmit={(e) => {
						e.preventDefault();
						form.handleSubmit();
					}}
					className="space-y-4"
				>
					<form.Field name="name">
						{(field) => (
							<div className="space-y-2">
								<Label htmlFor="name">Name</Label>
								<Input
									id="name"
									value={field.state.value}
									onChange={(e) => field.handleChange(e.target.value)}
									placeholder="Team name"
									maxLength={127}
								/>
							</div>
						)}
					</form.Field>

					<form.Field name="description">
						{(field) => (
							<div className="space-y-2">
								<Label htmlFor="description">Description</Label>
								<Textarea
									id="description"
									value={field.state.value}
									onChange={(e) => field.handleChange(e.target.value)}
									placeholder="What is this team about?"
									maxLength={500}
								/>
							</div>
						)}
					</form.Field>

					<form.Field name="max_members">
						{(field) => (
							<div className="space-y-2">
								<Label htmlFor="max_members">Max Members</Label>
								<Input
									id="max_members"
									type="number"
									value={field.state.value}
									onChange={(e) =>
										field.handleChange(Number(e.target.value))
									}
									min={1}
								/>
							</div>
						)}
					</form.Field>

					<Button
						type="submit"
						className="w-full"
						disabled={createTeam.isPending}
					>
						{createTeam.isPending ? "Creating..." : "Create Team"}
					</Button>
				</form>
			</DialogContent>
		</Dialog>
	);
}
