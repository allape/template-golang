import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  CrudyButton,
  Ellipsis,
  searchable,
} from "@allape/gocrud-react";
import { ICrudyButtonProps } from "@allape/gocrud-react/src/component/CrudyButton";
import { Form, Input, TableColumnsType } from "antd";
import { ReactElement, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { UserCrudy } from "../../api/user.ts";
import { IUser, IUserSearchParams } from "../../model/user.ts";

type IRecord = IUser;
type ISearchParams = IUserSearchParams;

export type IUserCrudyButtonProps = Partial<ICrudyButtonProps<IRecord>>;

export default function UserCrudyButton(
  props: IUserCrudyButtonProps,
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
        title: t("user.name"),
        dataIndex: "name",
        filtered: !!searchParams["like_name"],
        ...searchable(t("user.name"), (value) =>
          setSearchParams((old) => ({
            ...old,
            like_name: value,
          })),
        ),
      },
      {
        title: t("user.description"),
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
      name={t("user._")}
      columns={columns}
      crudy={UserCrudy}
      searchParams={searchParams}
      {...props}
    >
      <Form.Item
        name="name"
        label={t("user.name")}
        rules={[{ required: true }]}
      >
        <Input maxLength={200} placeholder={t("user.name")} />
      </Form.Item>
      <Form.Item name="description" label={t("user.description")}>
        <Input.TextArea
          maxLength={20000}
          rows={10}
          placeholder={t("user.description")}
        />
      </Form.Item>
    </CrudyButton>
  );
}
