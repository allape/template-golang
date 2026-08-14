import { CheckSquareOutlined, CloseSquareOutlined } from "@ant-design/icons";
import { Tag } from "antd";
import { ReactElement } from "react";

export interface IStatusTagProps {
  checked?: boolean;
}

export default function StatusTag({ checked }: IStatusTagProps): ReactElement {
  return checked ? (
    <Tag color="green">
      <CheckSquareOutlined />
    </Tag>
  ) : (
    <Tag color="red">
      <CloseSquareOutlined />
    </Tag>
  );
}
