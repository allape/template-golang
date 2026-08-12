import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  config,
  CrudyButton,
  EEEvent,
  Ellipsis,
  ICrudyButtonProps,
  NewCrudyButtonEventEmitter,
  Uploader,
} from "@allape/gocrud-react";
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
import { GalleryItemHandler } from "../../api/gallery.ts";
import { ItemCrudy, ItemTagHandler } from "../../api/item.ts";
import { IGallery } from "../../model/gallery.ts";
import { IItem, IItemSearchParams } from "../../model/item.ts";
import { ITag } from "../../model/tag.ts";
import GallerySelector from "../GallerySelector";
import TagSelector from "../TagSelector";

interface IItemModified extends IItem {
  _src?: string;
  _continuesUpload?: boolean;

  _galleryIds?: IGallery["id"][];
  _tagIds?: ITag["id"][];

  _tags?: ITag[];
  _galleries?: IGallery[];
}

type IRecord = IItemModified;
type ISearchParams = IItemSearchParams;

const crudy = ItemCrudy;

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

  const [searchParams] = useState<ISearchParams>(() => ({
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
      },
      {
        title: t("item.tags"),
        dataIndex: "_tags",
        render: (v: IRecord["_tags"]) =>
          v && v.length > 0
            ? v.map((i) => (
                <div key={i.id}>
                  <Tag color={i.color}>{i.name}</Tag>
                </div>
              ))
            : "---",
      },
      {
        title: t("item.galleries"),
        dataIndex: "_galleries",
        render: (v: IRecord["_galleries"]) =>
          v && v.length > 0
            ? v.map((i) => (
                <div key={i.id}>
                  <Tag>{i.name}</Tag>
                </div>
              ))
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
    [t],
  );

  const handleAfterListed = useCallback(
    async (records: IRecord[]): Promise<IRecord[]> => {
      if (records.length === 0) {
        return [];
      }

      records.forEach((r) => {
        r._src = `${config.SERVER_STATIC_URL}${r.src}`;
      });

      await ItemTagHandler.get<ITag, IItemModified>(
        "itemId",
        records,
        {},
        (item, tags) => {
          item._tags = tags;
          item._tagIds = tags.map((t) => t.id);
        },
      );

      await GalleryItemHandler.get<IGallery, IItemModified>(
        "itemId",
        records,
        {},
        (item, galleries) => {
          item._galleries = galleries;
          item._galleryIds = galleries.map((g) => g.id);
        },
      );

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
      const galleryIds: IGallery["id"][] =
        form.getFieldValue("_galleryIds") || [];
      await GalleryItemHandler.saveAfterDelete(
        "itemId",
        record.id,
        galleryIds.map((gi) => ({
          galleryId: gi,
          itemId: record.id,
        })),
      );

      const tagIds: ITag["id"][] = form.getFieldValue("_tagIds") || [];
      await ItemTagHandler.saveAfterDelete(
        "itemId",
        record.id,
        tagIds.map((ti) => ({
          itemId: record.id,
          tagId: ti,
        })),
      );

      const shouldStop: boolean = form.getFieldValue("_continuesUpload");

      form.resetFields();
      form.setFieldsValue({
        _continuesUpload: true,
        name: record.name,
      });

      if (shouldStop) {
        message.success(t("item.saved")).then();
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
    <CrudyButton<IRecord, ISearchParams>
      name={t("item._")}
      titleSearchField="like_name"
      columns={columns}
      crudy={crudy}
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
