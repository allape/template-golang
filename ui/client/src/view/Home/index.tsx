import { useLoading, useProxy } from "@allape/use-loading";
import { LoadingOutlined } from "@ant-design/icons";
import { AutoComplete, AutoCompleteProps, Empty, Input } from "antd";
import {
  KeyboardEvent,
  ReactElement,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import {
  getAllGalleries,
  IGalleryInfo,
  toImageURL,
} from "../../api/gallery.ts";
import styles from "./style.module.scss";

let lastSearchKeywords = "";

interface IGalleryInfoModified extends IGalleryInfo {
  _coverLink?: string;
}

export default function Home(): ReactElement {
  const navigate = useNavigate();

  const { t } = useTranslation();
  const { loading, execute } = useLoading();

  const searchTimerRef = useRef<number>(-1);

  const [keywords, keywordsRef, setKeywords] =
    useProxy<string>(lastSearchKeywords);

  const infoRef = useRef<IGalleryInfoModified[]>([]);
  const [galleries, setGalleries] = useState<IGalleryInfoModified[]>([]);

  const [autoCompleteOptions, setAutoCompleteOptions] = useState<
    AutoCompleteProps["options"]
  >([]);

  const handleSearch = useCallback(() => {
    clearTimeout(searchTimerRef.current);
    searchTimerRef.current = setTimeout(() => {
      lastSearchKeywords = keywordsRef.current;
      setGalleries(
        infoRef.current.filter((info) => {
          return (
            info.info.name.includes(keywordsRef.current) ||
            info.info.keywords.includes(keywordsRef.current) ||
            !!info.tags?.find((t) => t.name.includes(keywordsRef.current))
          );
        }),
      );
    }, 200) as unknown as number;
  }, [keywordsRef]);

  useEffect(() => {
    execute(async () => {
      const autoComOptions: AutoCompleteProps["options"] = [];

      const infos: IGalleryInfoModified[] = await getAllGalleries();
      infos.forEach((info) => {
        info._coverLink = toImageURL(info.info.id, 0, "thumbnail");

        info.tags?.map((tag) => {
          if (!autoComOptions.find((aco) => aco.value === tag.name)) {
            autoComOptions.push({
              value: tag.name,
              label: tag.name,
            });
          }
        });
      });

      infoRef.current = infos;
      setAutoCompleteOptions(autoComOptions);

      handleSearch();
    }).then();
  }, [execute, handleSearch]);

  useEffect(() => {
    return () => {
      clearTimeout(searchTimerRef.current);
    };
  }, []);

  const handleEnter = useCallback(
    (e: KeyboardEvent<HTMLInputElement>) => {
      handleSearch();
      (e.target as HTMLInputElement).blur();
    },
    [handleSearch],
  );

  return (
    <div className={styles.wrapper}>
      <div className={styles.search}>
        <AutoComplete options={autoCompleteOptions} showSearch>
          <Input
            className={styles.searchBox}
            type="search"
            size="large"
            placeholder={t("home.search")}
            value={keywords}
            onChange={(e) => {
              setKeywords(e.target.value);
              handleSearch();
            }}
            onPressEnter={handleEnter}
            onBlur={handleSearch}
            allowClear
          />
        </AutoComplete>
      </div>
      {loading ? (
        <div className={styles.loading}>
          <LoadingOutlined />
        </div>
      ) : undefined}
      <div className={styles.galleries}>
        {galleries.length === 0 ? (
          <div className={styles.empty}>
            <Empty />
          </div>
        ) : (
          galleries.map((gallery) => (
            <div
              key={gallery.info.id}
              data-type="gallery"
              className={styles.gallery}
              title={gallery.info.name}
              onClick={() => navigate(`/gallery/${gallery.info.id}`)}
              tabIndex={0}
            >
              <img
                className={styles.cover}
                loading="lazy"
                src={gallery._coverLink}
                alt={gallery.info.name}
              />
              <div className={styles.name}>{gallery.info.name}</div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
