"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@ui/components/button";
import { DataTable } from "@ui/components/ui/data-table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@ui/components/dropdown-menu";

import { Badge } from "@ui/components/ui/badge";

import { CalendarDots, NotePencil, CheckCircle } from "@ui/components/react-icons";
import { Service } from "@/web-sdk";
import { useAuth } from "@/components/auth-provider";

// Define the Blog type
type BlogPost = {
  organization_id: string;
  id: string;
  title: string;
  preview: string;
  status: "draft" | "ready" | "published"; // Assuming status has specific allowed values
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
  tags: string[];
  storage_path: string;
  created_by: string;
  updated_by: string;
};

export default function Page({ params }: { params: { organizationId: string } }) {
  const { user } = useAuth();
  const [blogs, setBlogs] = useState<BlogPost[]>();
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (!user) return;
    const getDocuments = async () => {
      try {
        const documents = await Service.makeAuthenticatedRequest("documents", "GET", null, {
          organization_id: params.organizationId,
        });
        setBlogs(documents);
      } catch (error) {
        console.log(error);
      }
    };

    getDocuments();
  }, [user]);

  const deleteBlog = (blogId: number) => {
    // setOpen(!open);
    // setBlogs((prevBlogs) => {
    //   const updatedBlogs = prevBlogs.filter((blog) => blog.id !== blogId);
    //   return updatedBlogs;
    // });
  };

  const tableData = () => {
    return blogs?.map((blog) => {
      // Generate tag badges
      const tagsBadges = blog.tags.map((tag) => (
        <Badge
          key={tag}
          className="bg-gray-200 text-black dark:bg-gray-700 dark:text-white"
          variant="outline"
        >
          {tag}
        </Badge>
      ));

      // Map status to icons
      const getStatusIcon = (status: string) => {
        switch (status) {
          case "Scheduled":
            return <CalendarDots className="w-4 h-4" />;
          case "Draft":
            return <NotePencil className="w-4 h-4" />;
          case "Published":
            return <CheckCircle className="w-4 h-4" />;
          default:
            return null;
        }
      };

      return {
        id: blog.id,
        title: blog.title,
        tags: (
          <div className="flex flex-wrap gap-2 overflow-x-auto whitespace-nowrap">{tagsBadges}</div>
        ),
        status: (
          <div className="flex items-center space-x-2">
            {getStatusIcon(blog.status)}
            <span>{blog.status}</span>
          </div>
        ),
        preview: (
          <div className="flex items-start mb-2 max-w-[70%] overflow-hidden">
            <p className="flex-1 ml-2 truncate">{blog.preview}</p>
          </div>
        ),
      };
    });
  };

  const columns = [
    {
      id: "title",
      accessorKey: "title",
      header: "Title",
      enableSorting: true,
    },
    {
      accessorKey: "tags",
      header: "Tags",
      enableSorting: false,
      cell: ({ getValue }) => <div>{getValue()}</div>,
    },
    {
      accessorKey: "preview",
      header: "Preview",
      enableSorting: false,
      cell: ({ getValue }) => <div>{getValue()}</div>,
    },
    {
      accessorKey: "status",
      header: "Status",
      enableSorting: true,
      cell: ({ getValue }) => <div>{getValue()}</div>,
    },
    {
      accessorKey: "actions",
      header: "Actions",
      id: "actions",
      enableSorting: false,
      cell: ({ row }) => {
        const blogId = row.original.id;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                ...
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>Edit</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => {
                  deleteBlog(blogId);
                }}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  return (
    <div className="space-y-4">
      {blogs?.length > 0 && (
        <DataTable columns={columns} data={tableData()} filterColumnId={"title"} />
      )}
    </div>
  );
}
