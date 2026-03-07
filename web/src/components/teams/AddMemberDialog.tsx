"use client";

import { useEffect, useRef, useState } from "react";
import { UserPlus, X } from "lucide-react";
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
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { useAddTeamMember, useTeamMembers } from "@/hooks/useTeams";
import { useSearchUsers } from "@/hooks/useSearchUsers";
import type { UserSearchResult } from "@/api/users";

export function AddMemberDialog({ teamId }: { teamId: string }) {
	const [open, setOpen] = useState(false);
	const [searchInput, setSearchInput] = useState("");
	const [debouncedQuery, setDebouncedQuery] = useState("");
	const [selectedUser, setSelectedUser] = useState<UserSearchResult | null>(
		null,
	);
	const [showDropdown, setShowDropdown] = useState(false);
	const [role, setRole] = useState<"member" | "admin">("member");
	const dropdownRef = useRef<HTMLDivElement>(null);

	const addMember = useAddTeamMember(teamId);
	const { data: members } = useTeamMembers(teamId);
	const { data: users, isLoading } = useSearchUsers(debouncedQuery);

	const memberIds = new Set(members?.map((m) => m.user_id));
	const filteredUsers = users?.filter((u) => !memberIds.has(u.id));

	// Debounce search input
	useEffect(() => {
		const timer = setTimeout(() => {
			setDebouncedQuery(searchInput);
		}, 300);
		return () => clearTimeout(timer);
	}, [searchInput]);

	// Show dropdown when we have results
	useEffect(() => {
		if (filteredUsers && filteredUsers.length > 0 && searchInput.length > 0) {
			setShowDropdown(true);
		} else {
			setShowDropdown(false);
		}
	}, [filteredUsers, searchInput]);

	// Close dropdown on outside click
	useEffect(() => {
		function handleClick(e: MouseEvent) {
			if (
				dropdownRef.current &&
				!dropdownRef.current.contains(e.target as Node)
			) {
				setShowDropdown(false);
			}
		}
		document.addEventListener("mousedown", handleClick);
		return () => document.removeEventListener("mousedown", handleClick);
	}, []);

	function handleSelectUser(user: UserSearchResult) {
		setSelectedUser(user);
		setSearchInput("");
		setDebouncedQuery("");
		setShowDropdown(false);
	}

	function handleClearUser() {
		setSelectedUser(null);
	}

	function resetState() {
		setSearchInput("");
		setDebouncedQuery("");
		setSelectedUser(null);
		setShowDropdown(false);
		setRole("member");
	}

	async function handleSubmit(e: React.FormEvent) {
		e.preventDefault();
		if (!selectedUser) {
			toast.error("Please select a user");
			return;
		}
		try {
			await addMember.mutateAsync({
				user_id: selectedUser.id,
				role,
			});
			toast.success("Member added");
			setOpen(false);
			resetState();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : "Failed to add member");
		}
	}

	return (
		<Dialog
			open={open}
			onOpenChange={(v) => {
				setOpen(v);
				if (!v) resetState();
			}}
		>
			<DialogTrigger asChild>
				<Button variant="outline">
					<UserPlus className="mr-2 size-4" />
					Add Member
				</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Add Team Member</DialogTitle>
				</DialogHeader>
				<form onSubmit={handleSubmit} className="space-y-4">
					<div className="space-y-2">
						<Label>User</Label>
						{selectedUser ? (
							<div className="flex items-center gap-2 rounded-md border bg-muted px-3 py-2">
								<span className="flex-1 text-sm">
									{selectedUser.name}{" "}
									<span className="text-muted-foreground">
										({selectedUser.email})
									</span>
								</span>
								<button
									type="button"
									onClick={handleClearUser}
									className="text-muted-foreground hover:text-foreground"
								>
									<X className="size-4" />
								</button>
							</div>
						) : (
							<div className="relative" ref={dropdownRef}>
								<Input
									value={searchInput}
									onChange={(e) => setSearchInput(e.target.value)}
									placeholder="Search by name..."
									autoComplete="off"
								/>
								{showDropdown && (
									<div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-md">
										{filteredUsers?.map((user) => (
											<button
												key={user.id}
												type="button"
												className="flex w-full flex-col px-3 py-2 text-left text-sm hover:bg-accent"
												onClick={() => handleSelectUser(user)}
											>
												<span className="font-medium">{user.name}</span>
												<span className="text-muted-foreground text-xs">
													{user.email}
												</span>
											</button>
										))}
									</div>
								)}
								{isLoading && searchInput.length > 0 && (
									<div className="absolute z-50 mt-1 w-full rounded-md border bg-popover p-3 text-center text-muted-foreground text-sm shadow-md">
										Searching...
									</div>
								)}
								{!isLoading &&
									debouncedQuery.length > 0 &&
									filteredUsers?.length === 0 && (
										<div className="absolute z-50 mt-1 w-full rounded-md border bg-popover p-3 text-center text-muted-foreground text-sm shadow-md">
											No users found
										</div>
									)}
							</div>
						)}
					</div>

					<div className="space-y-2">
						<Label>Role</Label>
						<Select
							value={role}
							onValueChange={(v) => setRole(v as "member" | "admin")}
						>
							<SelectTrigger>
								<SelectValue />
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="member">Member</SelectItem>
								<SelectItem value="admin">Admin</SelectItem>
							</SelectContent>
						</Select>
					</div>

					<Button
						type="submit"
						className="w-full"
						disabled={addMember.isPending || !selectedUser}
					>
						{addMember.isPending ? "Adding..." : "Add Member"}
					</Button>
				</form>
			</DialogContent>
		</Dialog>
	);
}
