import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Card } from "@astryxdesign/core/Card";
import { CodeBlock } from "@astryxdesign/core/CodeBlock";
import { Button } from "@astryxdesign/core/Button";
import { TextInput } from "@astryxdesign/core/TextInput";
import { Switch } from "@astryxdesign/core/Switch";
import { Banner } from "@astryxdesign/core/Banner";
import Shell from "./Shell";

// Mirrors internal/account's User (GET /api/account).
interface User {
  ID: number;
  Email: string;
  Username: string;
  PublicProfile: boolean;
  CreatedAt: string;
}

function AccountContent() {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [confirmUsername, setConfirmUsername] = useState("");
  const [deleteError, setDeleteError] = useState<string | null>(null);

  function load() {
    fetch("/api/account", { credentials: "include" })
      .then((res) => {
        if (res.status === 401) {
          window.location.href = "/login";
          return null;
        }
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then((data) => data && setUser(data))
      .catch((err) => setError(String(err)));
  }
  useEffect(load, []);

  async function handleIssueToken() {
    const res = await fetch("/api/account/token", { method: "POST", credentials: "include" });
    if (res.ok) setToken((await res.json()).token);
  }
  async function handleRevokeToken() {
    await fetch("/api/account/token/revoke", { method: "POST", credentials: "include" });
    setToken(null);
  }
  async function handleTogglePublicProfile(enabled: boolean) {
    await fetch("/api/account/public-profile", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ enabled }),
    });
    load();
  }
  async function handleLogout() {
    await fetch("/api/logout", { method: "POST", credentials: "include" });
    window.location.href = "/login";
  }
  async function handleDelete() {
    setDeleteError(null);
    const res = await fetch("/api/account/delete", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ confirm_username: confirmUsername }),
    });
    if (!res.ok) {
      const payload = await res.json().catch(() => ({ error: `server returned ${res.status}` }));
      setDeleteError(payload.error ?? "Could not delete account.");
      return;
    }
    window.location.href = "/login";
  }

  if (error) return <Text color="accent">Couldn't load your account: {error}</Text>;
  if (!user) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={6}>
      <Heading level={1}>Account</Heading>
      <Card>
        <VStack gap={2}>
          <Text>{user.Email}</Text>
          <Text type="supporting">@{user.Username}</Text>
        </VStack>
      </Card>

      <VStack gap={2}>
        <Heading level={3}>API token</Heading>
        <Text type="supporting">Used by `breaklab sync` to push/pull progress from this account.</Text>
        {token && (
          <Card variant="green">
            <VStack gap={2}>
              <Text>Copy this token now — it won't be shown again.</Text>
              <CodeBlock code={token} language="plaintext" hasLanguageLabel={false} />
            </VStack>
          </Card>
        )}
        <HStack gap={2}>
          <Button label="Issue new token" variant="secondary" clickAction={handleIssueToken} />
          <Button label="Revoke token" variant="secondary" clickAction={handleRevokeToken} />
        </HStack>
      </VStack>

      <Switch
        label="Public profile"
        description={`Show your roadmap progress at /u/${user.Username}`}
        value={user.PublicProfile}
        changeAction={handleTogglePublicProfile}
      />

      <Button label="Log out" variant="secondary" clickAction={handleLogout} />

      <VStack gap={2}>
        <Heading level={3}>Delete account</Heading>
        {deleteError && <Banner status="error" title={deleteError} container="card" />}
        <TextInput label={`Type "${user.Username}" to confirm`} value={confirmUsername} onChange={setConfirmUsername} />
        <Button label="Delete account" variant="destructive" clickAction={handleDelete} />
      </VStack>
    </VStack>
  );
}

export default function Account() {
  return (
    <Shell>
      <AccountContent />
    </Shell>
  );
}
