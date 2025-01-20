"use client";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@ui/components/tabs";
import { Button } from "@ui/components/button";
import { toast } from "@ui/components/use-toast";
import WebPageTheme from "@/components/web-page-theme";
import WebBlogTheme from "@/components/web-blog-theme";
import { Service } from "@/web-sdk";
import { RocketLaunch } from "@ui/components/react-icons";

export default function Page({ params }: { params: { organizationId: string } }) {
  const handlePublish = async () => {
    try {
      const response = await Service.makeAuthenticatedRequest(
        "documents/publish",
        "POST",
        {},
        {
          organization_id: params.organizationId,
        }
      );
      if (response) {
        toast({
          title: "Build started",
          description: "Your build has been started",
        });
      }
    } catch (error) {
      toast({
        title: "Failed to start build",
        description: error.message,
      });
      console.error(error);
    }
  };

  return (
    <Tabs defaultValue="web">
      <div className="flex justify-between border-b p-4 w-full">
        <TabsList className="grid w-full grid-cols-2 max-w-[200px]">
          <TabsTrigger value="web">Web</TabsTrigger>
          <TabsTrigger value="blog">Blog</TabsTrigger>
        </TabsList>
        <Button onClick={handlePublish} className="flex items-center gap-1">
          <RocketLaunch className="w-4 h-4" /> Deploy
        </Button>
      </div>

      <div>
        <TabsContent value="web">
          <WebPageTheme />
        </TabsContent>
        <TabsContent value="blog">
          <WebBlogTheme />
        </TabsContent>
      </div>
    </Tabs>
  );
}
