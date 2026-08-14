import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  config,
  CrudyButton,
  EEEvent,
  Ellipsis,
  ICrudyButtonProps,
  NewCrudyButtonEventEmitter,
} from "@allape/gocrud-react";
import {
  Avatar,
  Divider,
  Form,
  FormInstance,
  Input,
  InputNumber,
  TableColumnsType,
  Tag,
} from "antd";
import { IGallery } from "common/src/model/gallery.ts";
import { IItem } from "common/src/model/item.ts";
import { ITag } from "common/src/model/tag.ts";
import {
  ChangeEvent,
  ReactElement,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { useTranslation } from "react-i18next";
import { GalleryItemHandler } from "../../api/gallery.ts";
import { ItemCrudy, ItemTagHandler, upload } from "../../api/item.ts";
import { IItemSearchParams } from "../../model/item.ts";
import GallerySelector from "../GallerySelector";
import TagSelector from "../TagSelector";

interface IItemModified extends IItem {
  _src?: string;
  _thumbnail?: string;

  _file?: File;

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

  const fileRef = useRef<File>();

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
            src={record._thumbnail}
            shape="square"
            style={{ cursor: "pointer" }}
            size={100}
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
        r._thumbnail = `${config.SERVER_STATIC_URL}${r.thumbnail}`;
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

  const handleAfterSaved = useCallback(
    async (record: IRecord, form: FormInstance<IRecord>) => {
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
    },
    [],
  );

  const handleSave = useCallback(async (record: IRecord): Promise<IRecord> => {
    return upload(
      fileRef.current as File,
      {
        ...record,
        _galleryIds: undefined,
        _tagIds: undefined,
        _file: undefined,
      } as IRecord,
    );
  }, []);

  const handleFileChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    fileRef.current = e.target.files?.[0];
  }, []);

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
      afterSaved={handleAfterSaved}
      onSave={handleSave}
      emitter={emitter}
      {...props}
    >
      {(record?: Partial<IRecord>) => (
        <>
          <Form.Item
            name="_file"
            label={t("item.src")}
            rules={
              record?.id
                ? undefined
                : [{ required: true, message: t("item.srcRequired") }]
            }
            extra={record?.id ? t("item.fileExtra") : undefined}
          >
            <Input
              type="file"
              accept="image/*,video/*"
              onChange={handleFileChange}
            />
          </Form.Item>

          <Divider />

          <Form.Item name="_galleryIds" label={t("item.galleries")}>
            <GallerySelector mode="multiple" />
          </Form.Item>

          <Form.Item name="_tagIds" label={t("item.tags")}>
            <TagSelector mode="multiple" />
          </Form.Item>

          <Divider />

          {/*<Form.Item*/}
          {/*  name="src"*/}
          {/*  label={t("item.src")}*/}
          {/*  rules={[{ required: true, message: t("item.srcRequired") }]}*/}
          {/*>*/}
          {/*  <Uploader serverURL={config.SERVER_STATIC_URL} />*/}
          {/*</Form.Item>*/}

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
        </>
      )}
    </CrudyButton>
  );
}
