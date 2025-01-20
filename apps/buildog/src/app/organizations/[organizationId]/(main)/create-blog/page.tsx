"use client";

import { Service } from "@/web-sdk";
import { Input } from "@ui/components/input";
import { Textarea } from "@ui/components/textarea";
import React, { useState } from "react";

export default function Page({ params }) {
  const { organizationId } = params;
  const [markdownText, setMarkdownText] = useState<string>("");
  const [documentName, setDocumentName] = useState<string>("");
  const [status, setStatus] = useState<string>("");

  const handleUpload = async () => {
    if (!markdownText.trim()) {
      setStatus("Please enter some Markdown content.");
      return;
    }

    if (!documentName.trim()) {
      setStatus("Please enter a document name.");
      return;
    }

    const sanitizedFileName = documentName.replaceAll(" ", "-");
    const fileName = sanitizedFileName.endsWith(".md")
      ? sanitizedFileName
      : `${sanitizedFileName}.md`;

    const markdownBlob = new Blob([markdownText], { type: "text/markdown" });
    const markdownFile = new File([markdownBlob], fileName, {
      type: "text/markdown",
    });

    const formData = new FormData();
    formData.append("file", markdownFile);
    formData.append("title", documentName);
    formData.append("preview", "preview test string");
    formData.append("status", "draft");
    formData.append("tags", JSON.stringify(["react", "markdown", "test"]));

    try {
      setStatus("Uploading...");
      const response = await Service.makeAuthenticatedRequest("documents", "POST", formData, {
        organization_id: organizationId as string,
      });
      setStatus(`Upload successful: ${response.data}`);
    } catch (error: any) {
      setStatus(`Upload failed: ${error.response?.data || error.message}`);
    }
  };

  return (
    <div>
      <Input
        type="text"
        value={documentName}
        onChange={(e) => setDocumentName(e.target.value)}
        placeholder="Enter document name (e.g., myfile.md)"
      />
      <Textarea
        value={markdownText}
        onChange={(e) => setMarkdownText(e.target.value)}
        placeholder="Enter Markdown content here..."
      />
      <button onClick={handleUpload}>Upload</button>
      {status && <p>{status}</p>}
    </div>
  );
}
