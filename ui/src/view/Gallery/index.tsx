import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  CrudyTable,
  Ellipsis,
  searchable,
} from "@allape/gocrud-react";
import NewCrudyButtonEventEmitter from "@allape/gocrud-react/src/component/CrudyButton/eventemitter.ts";
import { Size } from "@allape/gocrud-react/src/hook/useMobile.ts";
import { UseLoadingReturn } from "@allape/use-loading/lib/hook/useLoading";
import { CameraOutlined, MoreOutlined } from "@ant-design/icons";
import {
  Button,
  Divider,
  Dropdown,
  Form,
  Input,
  InputNumber,
  MenuProps,
  Switch,
  TableColumnsType,
  Tag,
} from "antd";
import { ReactElement, ReactNode, useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { GalleryCrudy } from "../../api/gallery.ts";
import ItemCrudyButton from "../../component/ItemCrudyButton";
import TagCrudyButton from "../../component/TagCrudyButton";
import { IGallery, IGallerySearchParams } from "../../model/gallery.ts";
import { IItem, IItemSearchParams } from "../../model/item.ts";
import { ITag, ITagSearchParams } from "../../model/tag.ts";
import styles from './style.module.scss';

type IRecord = IGallery;
type ISearchParams = IGallerySearchParams;

const DefaultFormValue: Partial<IRecord> = {
  isPublic: false,
  priority: 0,
};

export default function Gallery(): ReactElement {
  const { t } = useTranslation();

  const emitter = useMemo(
    () => ({
      Tag: NewCrudyButtonEventEmitter<ITag, ITagSearchParams>(),
      Item: NewCrudyButtonEventEmitter<IItem, IItemSearchParams>(),
    }),
    [],
  );

  const [searchParams, setSearchParams] = useState<ISearchParams>(() => ({
    ...BaseSearchParams,
    orderByDefault: "1",
  }));

  const columns = useMemo<TableColumnsType<IRecord>>(
    () => [
      {
        title: t("id"),
        dataIndex: "id",
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
        filtered: !!searchParams["like_name"],
        ...searchable(t("gallery.name"), (value) =>
          setSearchParams((old) => ({
            ...old,
            like_name: value,
          })),
        ),
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

  return (
    <>
      <CrudyTable<IRecord>
        name={t("gallery._")}
        title={t("gallery._")}
        crudy={GalleryCrudy}
        columns={columns}
        searchParams={searchParams}
        defaultFormValue={DefaultFormValue}
        actions={handleActions}
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
            </div>
          </>
        }
      >
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
      </CrudyTable>
      <Dropdown className={styles.fixedMenu} menu={{ items: menus }}>
        <Button>
          <MoreOutlined />
        </Button>
      </Dropdown>
    </>
  );
}
