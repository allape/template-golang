import { BaseSearchParams } from "@allape/gocrud";
import {
  asDefaultPattern,
  CrudyTable,
  Ellipsis,
  searchable,
} from "@allape/gocrud-react";
import NewCrudyButtonEventEmitter from "@allape/gocrud-react/src/component/CrudyButton/eventemitter.ts";
import { MoreOutlined } from "@ant-design/icons";
import {
  Button,
  Divider,
  Dropdown,
  Form,
  Input,
  MenuProps,
  TableColumnsType,
} from "antd";
import { ReactElement, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { UserCrudy } from "../../api/user.ts";
import UserCrudyButton from "../../component/UserCrudyButton";
import UserSelector from "../../component/UserSelector";
import { IUser, IUserSearchParams } from "../../model/user.ts";
import styles from "./style.module.scss";

type IRecord = IUser;
type ISearchParams = IUserSearchParams;

export default function User(): ReactElement {
  const { t } = useTranslation();

  const UserCrudyEmitter = useMemo(
    () => NewCrudyButtonEventEmitter<IUser>(),
    [],
  );

  const [searchParams, setSearchParams] = useState<ISearchParams>(() => ({
    ...BaseSearchParams,
  }));

  const columns = useMemo<TableColumnsType<IRecord>>(
    () => [
      {
        title: t("id"),
        width: 50,
        dataIndex: "id",
        filtered: !!searchParams["in_id"],
        ...searchable<IRecord, IUser["id"]>(
          t("user.name"),
          (value) =>
            setSearchParams((old) => ({
              ...old,
              in_id: value ? [value] : undefined,
            })),
          (value, onChange) => (
            <UserSelector value={value} onChange={onChange} />
          ),
        ),
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

  const menus = useMemo<MenuProps["items"]>(
    () => [
      {
        key: "User",
        label: t("user._"),
        onClick: () => {
          UserCrudyEmitter.dispatchEvent("open");
        },
      },
    ],
    [UserCrudyEmitter, t],
  );

  return (
    <CrudyTable<IRecord>
      className={styles.wrapper}
      name={t("user._")}
      crudy={UserCrudy}
      columns={columns}
      searchParams={searchParams}
      titleExtra={
        <>
          <Divider type="vertical" />
          <div className={styles.windowed}>
            <UserCrudyButton emitter={UserCrudyEmitter} />
          </div>
          <div className={styles.mobile}>
            <Dropdown menu={{ items: menus }}>
              <Button>
                <MoreOutlined />
              </Button>
            </Dropdown>
          </div>
        </>
      }
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
    </CrudyTable>
  );
}
