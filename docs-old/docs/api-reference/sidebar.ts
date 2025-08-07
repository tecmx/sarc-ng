import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  apiSidebar: [
    {
      type: "doc",
      id: "sarc-ng-api",
      label: "API Overview"
    },
    {
      type: "category",
      label: "Buildings",
      link: {
        type: "doc",
        id: "buildings",
      },
      items: [
        {
          type: "doc",
          id: "get-all-buildings",
          label: "Get all buildings",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "create-building",
          label: "Create a new building",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "get-building-by-id",
          label: "Get building by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "update-building",
          label: "Update building",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "delete-building",
          label: "Delete building",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Classes",
      link: {
        type: "doc",
        id: "classes",
      },
      items: [
        {
          type: "doc",
          id: "get-all-classes",
          label: "Get all classes",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "create-class",
          label: "Create a new class",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "get-class-by-id",
          label: "Get class by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "update-class",
          label: "Update class",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "delete-class",
          label: "Delete class",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Lessons",
      link: {
        type: "doc",
        id: "lessons",
      },
      items: [
        {
          type: "doc",
          id: "get-all-lessons",
          label: "Get all lessons",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "create-lesson",
          label: "Create a new lesson",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "get-lesson-by-id",
          label: "Get lesson by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "update-lesson",
          label: "Update lesson",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "delete-lesson",
          label: "Delete lesson",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Reservations",
      link: {
        type: "doc",
        id: "reservations",
      },
      items: [
        {
          type: "doc",
          id: "get-all-reservations",
          label: "Get all reservations",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "create-reservation",
          label: "Create a new reservation",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "get-reservation-by-id",
          label: "Get reservation by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "update-reservation",
          label: "Update reservation",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "delete-reservation",
          label: "Delete reservation",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Resources",
      link: {
        type: "doc",
        id: "resources",
      },
      items: [
        {
          type: "doc",
          id: "get-all-resources",
          label: "Get all resources",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "create-resource",
          label: "Create a new resource",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "get-resource-by-id",
          label: "Get resource by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "update-resource",
          label: "Update resource",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "delete-resource",
          label: "Delete resource",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Authentication",
      link: {
        type: "doc",
        id: "auth",
      },
      items: [
        {
          type: "doc",
          id: "login",
          label: "User login",
          className: "api-method post",
        },
      ],
    },
  ],
};

export default sidebars;