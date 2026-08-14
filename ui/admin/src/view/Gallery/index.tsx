import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  CrudyTable,
  Ellipsis,
  NewCrudyButtonEventEmitter,
  searchable,
} from "@allape/gocrud-react";
import { Size } from "@allape/gocrud-react/src/hook/useMobile.ts";
import { UseLoadingReturn } from "@allape/use-loading/lib/hook/useLoading";
import { CameraOutlined, MoreOutlined } from "@ant-design/icons";
import {
  Button,
  Divider,
  Dropdown,
  Form,
  FormInstance,
  Input,
  InputNumber,
  MenuProps,
  Switch,
  TableColumnsType,
  TableProps,
  Tag,
} from "antd";
import { IGallery } from "common/src/model/gallery.ts";
import { IItem } from "common/src/model/item.ts";
import { ITag } from "common/src/model/tag.ts";
import { IUserGallery } from "common/src/model/user.ts";
import { ReactElement, ReactNode, useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { GalleryCrudy, GalleryTagHandler } from "../../api/gallery.ts";
import ItemCrudyButton from "../../component/ItemCrudyButton";
import StatusTag from "../../component/StatusTag";
import TagCrudyButton from "../../component/TagCrudyButton";
import TagSelector from "../../component/TagSelector";
import UserGalleryCrudyButton from "../../component/UserGalleryCrudyButton";
import { IGallerySearchParams } from "../../model/gallery.ts";
import { IItemSearchParams } from "../../model/item.ts";
import { ITagSearchParams } from "../../model/tag.ts";
import { IUserGallerySearchParams } from "../../model/user.ts";
import styles from "./style.module.scss";

const TableScroll: TableProps["scroll"] = {
  y: "calc(100vh - 200px)",
  x: true,
};

interface IGalleryModified extends IGallery {
  _tags?: ITag[];
  _tagIds?: ITag["id"][];
}

type IRecord = IGalleryModified;
type ISearchParams = IGallerySearchParams;

const crudy = GalleryCrudy;

const DefaultFormValue: Partial<IRecord> = {
  isPublic: false,
  priority: 0,
  enabled: true,
};

export default function Gallery(): ReactElement {
  const { t } = useTranslation();

  const emitter = useMemo(
    () => ({
      Tag: NewCrudyButtonEventEmitter<ITag, ITagSearchParams>(),
      Item: NewCrudyButtonEventEmitter<IItem, IItemSearchParams>(),
      UserGallery: NewCrudyButtonEventEmitter<
        IUserGallery,
        IUserGallerySearchParams
      >(),
    }),
    [],
  );

  const [searchParams, setSearchParams] = useState<ISearchParams>(() => ({
    ...BaseSearchParams,
    sortByPriorityThenUpdatedAt: true,
  }));

  const columns = useMemo<TableColumnsType<IRecord>>(
    () => [
      {
        title: t("id"),
        dataIndex: "id",
        render: (v: IRecord["id"], record: IRecord) => (
          <div>
            <StatusTag checked={record.enabled} /> {v}
          </div>
        ),
      },
      {
        title: t("gallery.tags"),
        dataIndex: "_tags",
        render: (tags?: ITag[]) =>
          tags?.map((tag) => (
            <div key={tag.id} style={{ marginBottom: "5px" }}>
              <Tag color={tag.color}>{tag.name}</Tag>
            </div>
          )) || "---",
      },
      {
        title: t("gallery.priority"),
        dataIndex: "priority",
      },
      {
        title: t("gallery.isPublic"),
        dataIndex: "isPublic",
        render: (v) =>
          v ? (
            <Tag color="red">{t("gallery.isPublicYesOrNo.yes")}</Tag>
          ) : (
            <Tag color="green">{t("gallery.isPublicYesOrNo.no")}</Tag>
          ),
      },
      {
        title: t("gallery.name"),
        dataIndex: "name",
      },
      {
        title: t("gallery.createdBy"),
        dataIndex: "createdBy",
        render: (v) => <Tag>{v}</Tag>,
        filtered: !!searchParams["createdBy"],
        ...searchable(t("gallery.createdBy"), (value) =>
          setSearchParams((old) => ({
            ...old,
            createdBy: value,
          })),
        ),
      },
      {
        title: t("gallery.keywords"),
        dataIndex: "keywords",
        render: (v) => <Ellipsis>{v}</Ellipsis>,
      },
      {
        title: t("gallery.description"),
        dataIndex: "description",
        render: (v) => <Ellipsis>{v}</Ellipsis>,
      },
      {
        title: t("createdAt"),
        dataIndex: "createdAt",
        render: asDefaultPattern,
      },
      {
        title: t("updatedAt"),
        dataIndex: "updatedAt",
        render: asDefaultPattern,
      },
    ],
    [searchParams, t],
  );

  const menus = useMemo<MenuProps["items"]>(
    () => [
      {
        key: "Tag",
        label: `${t("gocrud.manage")} ${t("tag._")}`,
        onClick: () => {
          emitter.Tag.dispatchEvent("open");
        },
      },
      {
        key: "Item",
        label: `${t("gocrud.manage")} ${t("item._")}`,
        onClick: () => {
          emitter.Item.dispatchEvent("open", {
            in_galleryId: undefined,
          });
        },
      },
      {
        key: "UserGallery",
        label: `${t("gocrud.manage")} ${t("user.gallery._")}`,
        onClick: () => {
          emitter.UserGallery.dispatchEvent("open");
        },
      },
    ],
    [emitter, t],
  );

  const handleActions = useCallback(
    (options: {
      record: IRecord;
      execute: UseLoadingReturn["execute"];
      size?: Size;
    }): ReactNode => {
      return (
        <>
          <Button
            type="link"
            onClick={() => {
              emitter.Item.dispatchEvent("open", {
                in_galleryId: [options.record.id],
              });
            }}
            title={t("gallery.photos")}
            size={options.size}
          >
            <CameraOutlined />
          </Button>
        </>
      );
    },
    [emitter, t],
  );

  const handleAfterListed = useCallback(
    async (records: IRecord[]): Promise<IRecord[]> => {
      await GalleryTagHandler.get<ITag, IGalleryModified>(
        "galleryId",
        records,
        {},
        (gallery, tags) => {
          gallery._tags = tags;
          gallery._tagIds = tags.map((t) => t.id);
        },
      );
      return records;
    },
    [],
  );

  const handleBeforeSave = useCallback((record: IRecord) => {
    delete record._tagIds;
  }, []);

  const handleAfterSaved = useCallback(
    async (record: IRecord, form: FormInstance<IRecord>) => {
      const tagIds: ITag["id"][] | undefined = form.getFieldValue("_tagIds");
      if (tagIds?.length) {
        await GalleryTagHandler.saveAfterDelete(
          "galleryId",
          record.id,
          tagIds.map((ti) => ({
            galleryId: record.id,
            tagId: ti,
          })),
        );
      }
    },
    [],
  );

  return (
    <>
      <CrudyTable<IRecord, ISearchParams>
        name={t("gallery._")}
        title={t("gallery._")}
        titleSearchField="keywords"
        crudy={crudy}
        columns={columns}
        searchParams={searchParams}
        defaultFormValue={DefaultFormValue}
        actions={handleActions}
        afterListed={handleAfterListed}
        beforeSave={handleBeforeSave}
        afterSaved={handleAfterSaved}
        scroll={TableScroll}
        titleExtra={
          <>
            <Divider type="vertical" />
            <Dropdown menu={{ items: menus }}>
              <Button>
                <MoreOutlined />
              </Button>
            </Dropdown>
            <div style={{ display: "none" }}>
              <TagCrudyButton emitter={emitter.Tag} />
              <ItemCrudyButton emitter={emitter.Item} />
              <UserGalleryCrudyButton emitter={emitter.UserGallery} />
            </div>
          </>
        }
      >
        <Form.Item name="_tagIds" label={t("gallery.tags")}>
          <TagSelector mode="multiple" />
        </Form.Item>

        <Form.Item name="isPublic" label={t("gallery.isPublic")}>
          <Switch
            checkedChildren={t("gallery.isPublicYesOrNo.yes")}
            unCheckedChildren={t("gallery.isPublicYesOrNo.no")}
          />
        </Form.Item>

        <Form.Item name="priority" label={t("gallery.priority")}>
          <InputNumber
            precision={0}
            step={1}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
            placeholder={t("gallery.priority")}
          />
        </Form.Item>

        <Form.Item
          name="name"
          label={t("gallery.name")}
          rules={[{ required: true }]}
        >
          <Input maxLength={200} placeholder={t("gallery.name")} />
        </Form.Item>

        <Form.Item name="keywords" label={t("gallery.keywords")}>
          <Input maxLength={500} placeholder={t("gallery.keywords")} />
        </Form.Item>

        <Form.Item name="createdBy" label={t("gallery.createdBy")}>
          <Input maxLength={200} placeholder={t("gallery.createdBy")} />
        </Form.Item>

        <Form.Item name="description" label={t("gallery.description")}>
          <Input.TextArea
            maxLength={20000}
            rows={10}
            placeholder={t("gallery.description")}
          />
        </Form.Item>

        <Form.Item name="enabled" label={t("gallery.enabled")}>
          <Switch
            checkedChildren={t("enabledYesOrNo.yes")}
            unCheckedChildren={t("enabledYesOrNo.no")}
          />
        </Form.Item>
      </CrudyTable>
      <Dropdown className={styles.fixedMenu} menu={{ items: menus }}>
        <Button>
          <MoreOutlined />
        </Button>
      </Dropdown>
    </>
  );
}
