import { useState } from "react";
import { Theme } from "@astryxdesign/core/theme";
import { Center } from "@astryxdesign/core/Center";
import { VStack } from "@astryxdesign/core/VStack";
import { Card } from "@astryxdesign/core/Card";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { TextInput } from "@astryxdesign/core/TextInput";
import { Button } from "@astryxdesign/core/Button";
import { Banner } from "@astryxdesign/core/Banner";
import { Link } from "@astryxdesign/core/Link";
import { themes, defaultThemeSlug } from "../astryx-themes";

// Standalone auth page: no Shell/nav, same as the logged-out login
// page — there's no account yet to show app chrome for.
function SignupForm() {
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSignup() {
    setError("");
    const res = await fetch("/api/signup", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, username, password }),
    });
    if (!res.ok) {
      const payload = await res.json().catch(() => ({ error: `server returned ${res.status}` }));
      setError(payload.error ?? "Could not create account.");
      return;
    }
    window.location.href = "/account";
  }

  return (
    <Center axis="both" padding={6}>
      <VStack gap={4} hAlign="center" width="100%" maxWidth={400}>
        <Card padding={8} width="100%">
          <VStack gap={4} hAlign="stretch">
            <VStack gap={1} hAlign="center">
              <Heading level={1}>Create your account</Heading>
              <Text type="body" color="secondary" size="sm">
                Sync progress across devices
              </Text>
            </VStack>
            {error && <Banner status="error" title={error} container="card" />}
            <TextInput label="Email" value={email} onChange={setEmail} type="email" size="lg" />
            <TextInput label="Username" value={username} onChange={setUsername} size="lg" />
            <TextInput label="Password" value={password} onChange={setPassword} type="password" size="lg" />
            <Button label="Sign up" variant="primary" size="lg" clickAction={handleSignup} />
            <Text type="supporting">
              Already have an account? <Link href="/login">Log in</Link>
            </Text>
          </VStack>
        </Card>
      </VStack>
    </Center>
  );
}

export default function Signup() {
  return (
    <Theme theme={themes[defaultThemeSlug]}>
      <SignupForm />
    </Theme>
  );
}
