import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  CrudyButton,
  Ellipsis,
  searchable,
} from "@allape/gocrud-react";
import { ICrudyButtonProps } from "@allape/gocrud-react/src/component/CrudyButton";
import { Form, Input, InputNumber, TableColumnsType } from "antd";
import { ReactElement, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { TagCrudy } from "../../api/tag.ts";
import { ITag, ITagSearchParams } from "../../model/tag.ts";

type IRecord = ITag;
type ISearchParams = ITagSearchParams;

const DefaultFormValue: Partial<IRecord> = {
  priority: 0,
};

export type ITagCrudyButtonProps = Partial<ICrudyButtonProps<IRecord>>;

export default function TagCrudyButton(
  props: ITagCrudyButtonProps,
): ReactElement {
  const { t } = useTranslation();

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
        title: t("tag.priority"),
        dataIndex: "priority",
      },
      {
        title: t("tag.name"),
        dataIndex: "name",
        filtered: !!searchParams["like_name"],
        ...searchable(t("tag.name"), (value) =>
          setSearchParams((old) => ({
            ...old,
            like_name: value,
          })),
        ),
      },
      {
        title: t("tag.alias"),
        dataIndex: "alias",
        filtered: !!searchParams["like_alias"],
        ...searchable(t("tag.alias"), (value) =>
          setSearchParams((old) => ({
            ...old,
            like_alias: value,
          })),
        ),
      },
      {
        title: t("tag.description"),
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

  return (
    <CrudyButton
      name={t("tag._")}
      columns={columns}
      crudy={TagCrudy}
      searchParams={searchParams}
      defaultFormValue={DefaultFormValue}
      {...props}
    >
      <Form.Item name="priority" label={t("tag.priority")}>
        <InputNumber
          precision={0}
          step={1}
          min={Number.MIN_SAFE_INTEGER}
          max={Number.MAX_SAFE_INTEGER}
          placeholder={t("tag.priority")}
        />
      </Form.Item>

      <Form.Item name="name" label={t("tag.name")} rules={[{ required: true }]}>
        <Input maxLength={50} placeholder={t("tag.name")} />
      </Form.Item>

      <Form.Item name="alias" label={t("tag.alias")}>
        <Input maxLength={200} placeholder={t("tag.alias")} />
      </Form.Item>

      <Form.Item name="description" label={t("tag.description")}>
        <Input.TextArea
          maxLength={20000}
          rows={10}
          placeholder={t("tag.description")}
        />
      </Form.Item>
    </CrudyButton>
  );
}
