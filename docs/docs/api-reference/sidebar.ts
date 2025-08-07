import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebar: SidebarsConfig = {
  apisidebar: [
    {
      type: "doc",
      id: "api-reference/sarc-ng-api-reference",
    },
    {
      type: "category",
      label: "buildings",
      link: {
        type: "doc",
        id: "api-reference/buildings",
      },
      items: [
        {
          type: "doc",
          id: "api-reference/get-all-buildings",
          label: "Get all buildings",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/create-building",
          label: "Create a new building",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api-reference/get-building-by-id",
          label: "Get building by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/update-building",
          label: "Update building",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "api-reference/delete-building",
          label: "Delete building",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "classes",
      link: {
        type: "doc",
        id: "api-reference/classes",
      },
      items: [
        {
          type: "doc",
          id: "api-reference/get-all-classes",
          label: "Get all classes",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/create-class",
          label: "Create a new class",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api-reference/get-class-by-id",
          label: "Get class by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/update-class",
          label: "Update class",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "api-reference/delete-class",
          label: "Delete class",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "lessons",
      link: {
        type: "doc",
        id: "api-reference/lessons",
      },
      items: [
        {
          type: "doc",
          id: "api-reference/get-all-lessons",
          label: "Get all lessons",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/create-lesson",
          label: "Create a new lesson",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api-reference/get-lesson-by-id",
          label: "Get lesson by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/update-lesson",
          label: "Update lesson",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "api-reference/delete-lesson",
          label: "Delete lesson",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "reservations",
      link: {
        type: "doc",
        id: "api-reference/reservations",
      },
      items: [
        {
          type: "doc",
          id: "api-reference/get-all-reservations",
          label: "Get all reservations",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/create-reservation",
          label: "Create a new reservation",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api-reference/get-reservation-by-id",
          label: "Get reservation by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/update-reservation",
          label: "Update reservation",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "api-reference/delete-reservation",
          label: "Delete reservation",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "resources",
      link: {
        type: "doc",
        id: "api-reference/resources",
      },
      items: [
        {
          type: "doc",
          id: "api-reference/get-all-resources",
          label: "Get all resources",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/create-resource",
          label: "Create a new resource",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "api-reference/get-resource-by-id",
          label: "Get resource by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "api-reference/update-resource",
          label: "Update resource",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "api-reference/delete-resource",
          label: "Delete resource",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "auth",
      link: {
        type: "doc",
        id: "api-reference/auth",
      },
      items: [
        {
          type: "doc",
          id: "api-reference/login",
          label: "User login",
          className: "api-method post",
        },
      ],
    },
  ],
};

export default sidebar.apisidebar;
