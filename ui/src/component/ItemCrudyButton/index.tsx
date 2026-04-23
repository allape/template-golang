import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  config,
  CrudyButton,
  EEEvent,
  Ellipsis,
  ICrudyButtonProps,
  searchable,
  Uploader,
} from "@allape/gocrud-react";
import NewCrudyButtonEventEmitter from "@allape/gocrud-react/src/component/CrudyButton/eventemitter.ts";
import {
  App,
  Avatar,
  Divider,
  Form,
  FormInstance,
  Input,
  InputNumber,
  Switch,
  TableColumnsType,
  Tag,
} from "antd";
import { ReactElement, useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  addItemToGalleries,
  GalleryCrudy,
  GalleryItemCrudy,
} from "../../api/gallery.ts";
import { addTagsToItem, ItemCrudy, ItemTagCrudy } from "../../api/item.ts";
import { TagCrudy } from "../../api/tag.ts";
import {
  IGallery,
  IGalleryItemSearchParams,
  IGallerySearchParams,
} from "../../model/gallery.ts";
import {
  IItem,
  IItemSearchParams,
  IItemTagSearchParams,
} from "../../model/item.ts";
import { ITag, ITagSearchParams } from "../../model/tag.ts";
import GallerySelector from "../GallerySelector";
import TagSelector from "../TagSelector";

interface ModifiedItem extends IItem {
  _src?: string;
  _continuesUpload?: boolean;

  _galleryIds?: IGallery["id"][];
  _tagIds?: ITag["id"][];

  _tags?: ITag[];
  _galleries?: IGallery[];
}

type IRecord = ModifiedItem;
type ISearchParams = IItemSearchParams;

const DefaultFormValue: Partial<IRecord> = {
  priority: 0,
};

export type IItemCrudyButtonProps = Partial<ICrudyButtonProps<IRecord>>;

export default function ItemCrudyButton({
  emitter = NewCrudyButtonEventEmitter<IItem, IItemSearchParams>(),
  ...props
}: IItemCrudyButtonProps): ReactElement {
  const { t } = useTranslation();
  const { message } = App.useApp();

  const [searchParams, setSearchParams] = useState<ISearchParams>(() => ({
    ...BaseSearchParams,
  }));

  const columns = useMemo<TableColumnsType<IRecord>>(
    () => [
      {
        title: t("id"),
        dataIndex: "id",
      },
      {
        title: t("item.priority"),
        dataIndex: "priority",
      },
      {
        title: t("item.src"),
        dataIndex: "src",
        render: (_, record) => (
          <Avatar
            src={record._src}
            style={{ cursor: "pointer" }}
            size={40}
            onClick={() => window.open(record._src, "_blank")}
          />
        ),
      },
      {
        title: t("item.name"),
        dataIndex: "name",
        render: (v) => v || "---",
        filtered: !!searchParams["like_name"],
        ...searchable(t("item.name"), (value) =>
          setSearchParams((old) => ({
            ...old,
            like_name: value,
          })),
        ),
      },
      {
        title: t("item.tags"),
        dataIndex: "_tags",
        render: (v: IRecord["_tags"]) =>
          v && v.length > 0
            ? v.map((i) => (
                <Tag key={i.id} color={i.color}>
                  {i.name}
                </Tag>
              ))
            : "---",
      },
      {
        title: t("item.galleries"),
        dataIndex: "_galleries",
        render: (v: IRecord["_galleries"]) =>
          v && v.length > 0
            ? v.map((i) => <Tag key={i.id}>{i.name}</Tag>)
            : "---",
      },
      {
        title: t("item.description"),
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

  const handleAfterListed = useCallback(
    async (records: IRecord[]): Promise<IRecord[]> => {
      if (records.length === 0) {
        return [];
      }

      records.forEach((r) => {
        r._src = `${config.SERVER_STATIC_URL}${r.src}`;
      });

      const itemIds = Array.from(new Set(records.map((r) => r.id)));

      const itemTags = await ItemTagCrudy.all<IItemTagSearchParams>({
        in_itemId: itemIds,
      });
      const tagIds = Array.from(new Set(itemTags.map((i) => i.tagId)));
      const tags =
        tagIds.length > 0
          ? await TagCrudy.all<ITagSearchParams>({
              in_id: tagIds,
            })
          : [];

      const galleryItems = await GalleryItemCrudy.all<IGalleryItemSearchParams>(
        {
          in_itemId: itemIds,
        },
      );
      const galleryIds = Array.from(
        new Set(galleryItems.map((i) => i.galleryId)),
      );
      const galleries =
        galleryIds.length > 0
          ? await GalleryCrudy.all<IGallerySearchParams>({
              in_id: galleryIds,
            })
          : [];

      records.forEach((i) => {
        const its = itemTags.filter((it) => it.itemId === i.id);
        i._tagIds = its.map((it) => it.tagId);
        i._tags = its
          .map((it) => tags.find((t) => t.id === it.tagId))
          .filter(Boolean) as ITag[];

        const gis = galleryItems.filter((gi) => gi.itemId === i.id);
        i._galleryIds = gis.map((gi) => gi.galleryId);
        i._galleries = gis
          .map((gi) => galleries.find((g) => g.id === gi.galleryId))
          .filter(Boolean) as IGallery[];
      });

      return records;
    },
    [],
  );

  const handleBeforeSave = useCallback((record: IRecord): IRecord => {
    delete record._continuesUpload;
    delete record._galleryIds;
    delete record._tagIds;
    return record;
  }, []);

  const handleAfterSaved = useCallback(
    async (record: IRecord, form: FormInstance<IRecord>): Promise<boolean> => {
      const galleryIds: number[] | undefined =
        form.getFieldValue("_galleryIds");
      if (galleryIds && galleryIds.length > 0) {
        await addItemToGalleries(record.id, galleryIds);
      }

      const tagIds: ITag["id"][] | undefined = form.getFieldValue("_tagIds");
      if (tagIds && tagIds.length > 0) {
        await addTagsToItem(tagIds, record.id);
      }

      const shouldStop: boolean = form.getFieldValue("_continuesUpload");

      form.resetFields();
      form.setFieldsValue({
        _continuesUpload: true,
        name: record.name,
      });

      if (shouldStop) {
        message.success(t("item.saved"));
        return false;
      }
      return true;
    },
    [message, t],
  );

  const [defaultFormValue, setDefaultFormValue] = useState<Partial<IRecord>>(
    () => DefaultFormValue,
  );

  useEffect(() => {
    const handleOpen = (e: EEEvent<"open", ISearchParams | undefined>) => {
      setDefaultFormValue((o) => ({
        ...o,
        _galleryIds: e.value?.in_galleryId,
      }));
    };
    emitter.addEventListener("open", handleOpen);
    return () => {
      emitter.removeEventListener("open", handleOpen);
    };
  }, [emitter]);

  return (
    <CrudyButton
      name={t("item._")}
      columns={columns}
      crudy={ItemCrudy}
      searchParams={searchParams}
      defaultFormValue={defaultFormValue}
      afterListed={handleAfterListed}
      beforeSave={handleBeforeSave}
      afterSaved={handleAfterSaved}
      emitter={emitter}
      {...props}
    >
      <Form.Item name="_continuesUpload" label={t("item.continuesUpload")}>
        <Switch />
      </Form.Item>

      <Divider />

      <Form.Item name="_galleryIds" label={t("item.galleries")}>
        <GallerySelector mode="multiple" />
      </Form.Item>

      <Form.Item name="_tagIds" label={t("item.tags")}>
        <TagSelector mode="multiple" />
      </Form.Item>

      <Divider />

      <Form.Item
        name="src"
        label={t("item.src")}
        rules={[{ required: true, message: t("item.srcRequired") }]}
      >
        <Uploader serverURL={config.SERVER_STATIC_URL} />
      </Form.Item>

      <Form.Item name="priority" label={t("item.priority")}>
        <InputNumber
          precision={0}
          step={1}
          min={Number.MIN_SAFE_INTEGER}
          max={Number.MAX_SAFE_INTEGER}
          placeholder={t("item.priority")}
        />
      </Form.Item>

      <Form.Item name="name" label={t("item.name")}>
        <Input maxLength={50} placeholder={t("item.name")} />
      </Form.Item>

      <Form.Item name="description" label={t("item.description")}>
        <Input.TextArea
          maxLength={20000}
          rows={10}
          placeholder={t("item.description")}
        />
      </Form.Item>
    </CrudyButton>
  );
}
