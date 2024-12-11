package diff

import (
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestDiffJson(t *testing.T) {
	jsonString, err := DiffJsonString(input1, input2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(jsonString, true))
}

var (
	input1 = `{
  "base": {
    "id": 144,
    "name": "AnchorView",
    "enum_id": 10,
    "desc": "poi锚点消费侧场景",
    "owners": [
      "fanhaodong.516",
      "yangkaihong.2022",
      "zhangyue.x",
      "wangjingpei"
    ],
    "resp_groups": [
      {
        "id": 4597,
        "field_group": {
          "id": 2337,
          "name": "AggregatedSku",
          "desc": "Sku聚合信息",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4609,
        "field_group": {
          "id": 2400,
          "name": "AtmosphereInfo",
          "desc": "资源位（废弃）",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4584,
        "field_group": {
          "id": 2427,
          "name": "BuyEntry",
          "desc": "商品购买链接",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4587,
        "field_group": {
          "id": 2407,
          "name": "ClueInfo",
          "desc": "商品线索信息",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4588,
        "field_group": {
          "id": 2401,
          "name": "ClueTrackParams",
          "desc": "线索埋点信息",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4600,
        "field_group": {
          "id": 2288,
          "name": "GovernanceInfo",
          "desc": "治理信息",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4602,
        "field_group": {
          "id": 2298,
          "name": "IniosMiniOdrBlackTag",
          "desc": "小程序黑名单",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4604,
        "field_group": {
          "id": 2447,
          "name": "InnerIndustryExtData",
          "desc": "行业内数据",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2252,
            "name": "product_biz_data_4_trip",
            "desc": "酒旅自己生产的商品业务数据，如GMV等",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4612,
        "field_group": {
          "id": 2393,
          "name": "IsMember",
          "desc": "商品是否是会员商品",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2022,
            "name": "member_buy_info",
            "desc": "会员品购买条件信息",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4599,
        "field_group": {
          "id": 2418,
          "name": "MallPoiInfo",
          "desc": "商场POI适用的门店信息(已废弃，请选用MallPoiInfoV2)",
          "ctime": 1729063128,
          "mtime": 1732243668
        },
        "extra": [
          {
            "id": 2130,
            "name": "product_mall_poi_detail",
            "desc": "(商场适用）商品的门店详情",
            "ctime": 1729063124,
            "mtime": 1729063124
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4601,
        "field_group": {
          "id": 2410,
          "name": "MarketingData",
          "desc": "营销活动字段。包括会员价、券包、秒杀等营销活动信息 ",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2187,
            "name": "product_promotion_plan",
            "desc": "商品关联的促销计划",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4607,
        "field_group": {
          "id": 2402,
          "name": "MarketingDePriceInfo",
          "desc": "营销价格信息（废弃）",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4611,
        "field_group": {
          "id": 2384,
          "name": "MarketingExpression",
          "desc": "营销表达信息；展示前端营销氛围",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4608,
        "field_group": {
          "id": 2396,
          "name": "MarketingPriceInfo",
          "desc": "商品价格，结构化数据；营销优惠券",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2266,
            "name": "group_user_tag",
            "desc": "商城用户身份",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2180,
            "name": "sub_item_info",
            "desc": "日历房的间夜信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2170,
            "name": "trip_2c_product_map",
            "desc": "酒旅组品2C品商品信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2181,
            "name": "trip_combo_sku_price",
            "desc": "酒旅组品日历房售卖价",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2002,
            "name": "trip_package_product_sub_product_info",
            "desc": "酒旅组品的子品",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4594,
        "field_group": {
          "id": 2356,
          "name": "NearestPoiInfo",
          "desc": "商品关联POI信息(已废弃,请选用NearestPoiInfoV2)",
          "ctime": 1729063127,
          "mtime": 1732243446
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4603,
        "field_group": {
          "id": 2299,
          "name": "OutsideExtData",
          "desc": "外卖行业信息（废弃）",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "extra": [
          {
            "id": 2023,
            "name": "outside_balance_coupon_shark",
            "desc": "余额膨胀风控",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4592,
        "field_group": {
          "id": 2467,
          "name": "ProductAttr",
          "desc": "商品属性",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2008,
            "name": "poi_platform_config",
            "desc": "当前poi的平台配置，POI供货商",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2019,
            "name": "product_consumer_marketing",
            "desc": "商品C端营销信息",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2248,
            "name": "product_ele_industry_data",
            "desc": "指定商品的饿了么行业信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2279,
            "name": "trip_shelf_config",
            "desc": "酒旅货架干预能力配置",
            "ctime": 1729063126,
            "mtime": 1729063126
          }
        ],
        "ctime": 1729063136,
        "mtime": 1733369694
      },
      {
        "id": 4591,
        "field_group": {
          "id": 2340,
          "name": "ProductBaseInfo",
          "desc": "商品基础信息，名称、类型、图片、售卖时间等(不再包含商品图片，如需使用图片请添加ProductImages)",
          "ctime": 1729063126,
          "mtime": 1732244262
        },
        "extra": [
          {
            "id": 2007,
            "name": "poi_service_config",
            "desc": "当前poi的服务类型配置",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2039,
            "name": "product_stock_agg_data",
            "desc": "商品维度库存库存信息",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2104,
            "name": "product_stock_agg_data_with_date",
            "desc": "商品维度库存库存信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4593,
        "field_group": {
          "id": 2387,
          "name": "ProductBaseInfo.Desc",
          "desc": "商品描述信息",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4598,
        "field_group": {
          "id": 2336,
          "name": "ProductBizAttr",
          "desc": "commerce业务属性",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4595,
        "field_group": {
          "id": 2313,
          "name": "ProductEntry",
          "desc": "商品详情链接",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "extra": [
          {
            "id": 2063,
            "name": "combo_product_scenic_shelf_come_into_calendar_room_end_date",
            "desc": "日历房上架结束日期",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2062,
            "name": "combo_product_scenic_shelf_come_into_calendar_room_start_date",
            "desc": "日历房上架开始日期",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2314,
            "name": "expect_price_type",
            "desc": "主推价",
            "ctime": 1729063126,
            "mtime": 1729063126
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 6969,
        "field_group": {
          "id": 2478,
          "name": "ProductImages",
          "desc": "-",
          "ctime": 1729482493,
          "mtime": 1729482493
        },
        "ctime": 1729588842,
        "mtime": 1729588842
      },
      {
        "id": 4610,
        "field_group": {
          "id": 2456,
          "name": "ProductMemberInfo",
          "desc": "是否是会员专属商品，会员与会员价相关埋点",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4590,
        "field_group": {
          "id": 2434,
          "name": "ProductPrice",
          "desc": "商品价格(废弃)，后续使用营销价格字段（XX）。包括单价、优惠数量、次卡的单次价格、价格后缀",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4583,
        "field_group": {
          "id": 2468,
          "name": "ProductSaleStatus",
          "desc": "商品的购买信息",
          "ctime": 1729063129,
          "mtime": 1729063129
        },
        "extra": [
          {
            "id": 2229,
            "name": "poi_delivery_info",
            "desc": "指定门店的配送信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4606,
        "field_group": {
          "id": 2334,
          "name": "RespMarketingIndustry",
          "desc": "行业营销信息（废弃）",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4589,
        "field_group": {
          "id": 2469,
          "name": "SaleVolume",
          "desc": "商品售卖统计信息，包括 7天/30天销量等",
          "ctime": 1729063129,
          "mtime": 1729063129
        },
        "extra": [
          {
            "id": 2014,
            "name": "ele_aggr_product_sold_cnt",
            "desc": "饿了么聚合品销量",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4585,
        "field_group": {
          "id": 2392,
          "name": "ServingPrivilege",
          "desc": "服务权益",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4605,
        "field_group": {
          "id": 2446,
          "name": "ShelfNotFilter",
          "desc": "忽略货架的过滤器",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2252,
            "name": "product_biz_data_4_trip",
            "desc": "酒旅自己生产的商品业务数据，如GMV等",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4586,
        "field_group": {
          "id": 2428,
          "name": "TagGroup",
          "desc": "商品标签",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2182,
            "name": "aggregated_product_info",
            "desc": "聚合品主品商品信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2252,
            "name": "product_biz_data_4_trip",
            "desc": "酒旅自己生产的商品业务数据，如GMV等",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2034,
            "name": "product_sku_data",
            "desc": "sku数据",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2176,
            "name": "serving_privilege",
            "desc": "商品权益信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          }
        ],
        "ctime": 1729063136,
        "mtime": 1733301373
      },
      {
        "id": 7990,
        "field_group": {
          "id": 2488,
          "name": "TravelCalendarInfoV2",
          "desc": "TravelCalendarInfoV2",
          "creator": "wangjingpei",
          "ctime": 1733301373,
          "mtime": 1733301373
        },
        "ctime": 1733301373,
        "mtime": -62135596800
      },
      {
        "id": 7991,
        "field_group": {
          "id": 2489,
          "name": "TripGrouponDetailV2",
          "desc": "TripGrouponDetailV2",
          "creator": "wangjingpei",
          "ctime": 1733301373,
          "mtime": 1733301373
        },
        "ctime": 1733301373,
        "mtime": -62135596800
      },
      {
        "id": 4596,
        "field_group": {
          "id": 2349,
          "name": "UserDecisionInfo",
          "desc": "用户的决策信息：点击、消费、是否收藏",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2263,
            "name": "product_order_info",
            "desc": "商品关联的订单信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2250,
            "name": "product_user_decision",
            "desc": "商品用户决策信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2262,
            "name": "product_user_info",
            "desc": "商品关联的用户(点赞、购买)信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      }
    ],
    "filters": [
      {
        "id": 1880,
        "name": "FilterByDiscountRevert",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 10,
        "psm": "data.life.commerce"
      },
      {
        "id": 1888,
        "name": "FilterByLibraExperiment",
        "desc": "根据实验信息过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 54000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1889,
        "name": "FilterByMicroAppInvisibility",
        "desc": "根据微应用不可见性过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 45000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1882,
        "name": "FilterByNotIndependentSale",
        "desc": "根据非独立销售信息过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 50000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1883,
        "name": "FilterByOutOfDate",
        "desc": "根据过期信息过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 47000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1891,
        "name": "FilterByProductStatus",
        "desc": "根据商品状态过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 46000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1892,
        "name": "FilterByProductStock",
        "desc": "根据商品库存过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 71000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1886,
        "name": "FilterBySaleChannel",
        "desc": "根据销售渠道过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 56000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1887,
        "name": "FilterBySpuStatus",
        "desc": "根据SPU状态过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 43000,
        "psm": "data.life.commerce"
      }
    ],
    "ctime": 1729063136,
    "mtime": 1729611223,
    "psm": "data.life.commerce"
  }
}
`
	input2 = `{
  "base": {
    "id": 144,
    "name": "AnchorView",
    "enum_id": 10,
    "desc": "poi锚点消费侧场景",
    "owners": [
      "fanhaodong.516",
      "yangkaihong.2022",
      "zhangyue.x",
      "wangjingpei"
    ],
    "resp_groups": [
      {
        "id": 4597,
        "field_group": {
          "id": 2337,
          "name": "AggregatedSku",
          "desc": "Sku聚合信息",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4609,
        "field_group": {
          "id": 2400,
          "name": "AtmosphereInfo",
          "desc": "资源位（废弃）",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4584,
        "field_group": {
          "id": 2427,
          "name": "BuyEntry",
          "desc": "商品购买链接",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4587,
        "field_group": {
          "id": 2407,
          "name": "ClueInfo",
          "desc": "商品线索信息",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4588,
        "field_group": {
          "id": 2401,
          "name": "ClueTrackParams",
          "desc": "线索埋点信息",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4600,
        "field_group": {
          "id": 2288,
          "name": "GovernanceInfo",
          "desc": "治理信息",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4602,
        "field_group": {
          "id": 2298,
          "name": "IniosMiniOdrBlackTag",
          "desc": "小程序黑名单",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4604,
        "field_group": {
          "id": 2447,
          "name": "InnerIndustryExtData",
          "desc": "行业内数据",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2252,
            "name": "product_biz_data_4_trip",
            "desc": "酒旅自己生产的商品业务数据，如GMV等",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4612,
        "field_group": {
          "id": 2393,
          "name": "IsMember",
          "desc": "商品是否是会员商品",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2022,
            "name": "member_buy_info",
            "desc": "会员品购买条件信息",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4599,
        "field_group": {
          "id": 2418,
          "name": "MallPoiInfo",
          "desc": "商场POI适用的门店信息(已废弃，请选用MallPoiInfoV2)",
          "ctime": 1729063128,
          "mtime": 1732243668
        },
        "extra": [
          {
            "id": 2130,
            "name": "product_mall_poi_detail",
            "desc": "(商场适用）商品的门店详情",
            "ctime": 1729063124,
            "mtime": 1729063124
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4601,
        "field_group": {
          "id": 2410,
          "name": "MarketingData",
          "desc": "营销活动字段。包括会员价、券包、秒杀等营销活动信息 ",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2187,
            "name": "product_promotion_plan",
            "desc": "商品关联的促销计划",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4607,
        "field_group": {
          "id": 2402,
          "name": "MarketingDePriceInfo",
          "desc": "营销价格信息（废弃）",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4611,
        "field_group": {
          "id": 2384,
          "name": "MarketingExpression",
          "desc": "营销表达信息；展示前端营销氛围",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4608,
        "field_group": {
          "id": 2396,
          "name": "MarketingPriceInfo",
          "desc": "商品价格，结构化数据；营销优惠券",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2266,
            "name": "group_user_tag",
            "desc": "商城用户身份",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2180,
            "name": "sub_item_info",
            "desc": "日历房的间夜信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2170,
            "name": "trip_2c_product_map",
            "desc": "酒旅组品2C品商品信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2181,
            "name": "trip_combo_sku_price",
            "desc": "酒旅组品日历房售卖价",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2002,
            "name": "trip_package_product_sub_product_info",
            "desc": "酒旅组品的子品",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4594,
        "field_group": {
          "id": 2356,
          "name": "111",
          "desc": "商品关联POI信息(已废弃,请选用NearestPoiInfoV2)",
          "ctime": 1729063127,
          "mtime": 1732243446
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4603,
        "field_group": {
          "id": 2299,
          "name": "OutsideExtData",
          "desc": "外卖行业信息（废弃）",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "extra": [
          {
            "id": 2023,
            "name": "outside_balance_coupon_shark",
            "desc": "余额膨胀风控",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4592,
        "field_group": {
          "id": 2467,
          "name": "ProductAttr",
          "desc": "商品属性",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2008,
            "name": "poi_platform_config",
            "desc": "当前poi的平台配置，POI供货商",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2019,
            "name": "product_consumer_marketing",
            "desc": "商品C端营销信息",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2248,
            "name": "product_ele_industry_data",
            "desc": "指定商品的饿了么行业信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2279,
            "name": "trip_shelf_config",
            "desc": "酒旅货架干预能力配置",
            "ctime": 1729063126,
            "mtime": 1729063126
          }
        ],
        "ctime": 1729063136,
        "mtime": 1733369694
      },
      {
        "id": 4591,
        "field_group": {
          "id": 2340,
          "name": "ProductBaseInfo",
          "desc": "商品基础信息，名称、类型、图片、售卖时间等(不再包含商品图片，如需使用图片请添加ProductImages)",
          "ctime": 1729063126,
          "mtime": 1732244262
        },
        "extra": [
          {
            "id": 2007,
            "name": "poi_service_config",
            "desc": "当前poi的服务类型配置",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2039,
            "name": "product_stock_agg_data",
            "desc": "商品维度库存库存信息",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2104,
            "name": "product_stock_agg_data_with_date",
            "desc": "商品维度库存库存信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4593,
        "field_group": {
          "id": 2387,
          "name": "ProductBaseInfo.Desc",
          "desc": "商品描述信息",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4598,
        "field_group": {
          "id": 2336,
          "name": "ProductBizAttr",
          "desc": "commerce业务属性",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4595,
        "field_group": {
          "id": 2313,
          "name": "ProductEntry",
          "desc": "商品详情链接",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "extra": [
          {
            "id": 2063,
            "name": "combo_product_scenic_shelf_come_into_calendar_room_end_date",
            "desc": "日历房上架结束日期",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2062,
            "name": "combo_product_scenic_shelf_come_into_calendar_room_start_date",
            "desc": "日历房上架开始日期",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2314,
            "name": "expect_price_type",
            "desc": "主推价",
            "ctime": 1729063126,
            "mtime": 1729063126
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 6969,
        "field_group": {
          "id": 2478,
          "name": "ProductImages",
          "desc": "-",
          "ctime": 1729482493,
          "mtime": 222
        },
        "ctime": 1729588842,
        "mtime": 1729588842
      },
      {
        "id": 22,
        "field_group": {
          "id": 2456,
          "name": "ProductMemberInfo",
          "desc": "是否是会员专属商品，会员与会员价相关埋点",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4590,
        "field_group": {
          "id": 2434,
          "name": "ProductPrice",
          "desc": "商品价格(废弃)，后续使用营销价格字段（XX）。包括单价、优惠数量、次卡的单次价格、价格后缀",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4583,
        "field_group": {
          "id": 2468,
          "name": "ProductSaleStatus",
          "desc": "商品的购买信息",
          "ctime": 1729063129,
          "mtime": 1729063129
        },
        "extra": [
          {
            "id": 2229,
            "name": "poi_delivery_info",
            "desc": "指定门店的配送信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4606,
        "field_group": {
          "id": 2334,
          "name": "RespMarketingIndustry",
          "desc": "行业营销信息（废弃）",
          "ctime": 1729063126,
          "mtime": 1729063126
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4589,
        "field_group": {
          "id": 2469,
          "name": "SaleVolume",
          "desc": "商品售卖统计信息，包括 7天/30天销量等",
          "ctime": 1729063129,
          "mtime": 1729063129
        },
        "extra": [
          {
            "id": 2014,
            "name": "ele_aggr_product_sold_cnt",
            "desc": "饿了么聚合品销量",
            "ctime": 1729063123,
            "mtime": 1729063123
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4585,
        "field_group": {
          "id": 2392,
          "name": "ServingPrivilege",
          "desc": "服务权益",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4605,
        "field_group": {
          "id": 2446,
          "name": "ShelfNotFilter",
          "desc": "忽略货架的过滤器",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2252,
            "name": "product_biz_data_4_trip",
            "desc": "酒旅自己生产的商品业务数据，如GMV等",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      },
      {
        "id": 4586,
        "field_group": {
          "id": 2428,
          "name": "TagGroup",
          "desc": "商品标签",
          "ctime": 1729063128,
          "mtime": 1729063128
        },
        "extra": [
          {
            "id": 2182,
            "name": "aggregated_product_info",
            "desc": "聚合品主品商品信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          },
          {
            "id": 2252,
            "name": "product_biz_data_4_trip",
            "desc": "酒旅自己生产的商品业务数据，如GMV等",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2020,
            "name": "product_marketing_data",
            "desc": "营销算价原始上下文：到手价，引导价，会员价格，秒杀，各种券和减优惠",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2034,
            "name": "product_sku_data",
            "desc": "sku数据",
            "ctime": 1729063123,
            "mtime": 1729063123
          },
          {
            "id": 2176,
            "name": "serving_privilege",
            "desc": "商品权益信息",
            "ctime": 1729063124,
            "mtime": 1729063124
          }
        ],
        "ctime": 1729063136,
        "mtime": 1733301373
      },
      {
        "id": 7990,
        "field_group": {
          "id": 2488,
          "name": "TravelCalendarInfoV2",
          "desc": "TravelCalendarInfoV2",
          "creator": "wangjingpei",
          "ctime": 1733301373,
          "mtime": 1733301373
        },
        "ctime": 1733301373,
        "mtime": -62135596800
      },
      {
        "id": 7991,
        "field_group": {
          "id": 2489,
          "name": "TripGrouponDetailV2",
          "desc": "TripGrouponDetailV2",
          "creator": "wangjingpei",
          "ctime": 1733301373,
          "mtime": 1733301373
        },
        "ctime": 1733301373,
        "mtime": -62135596800
      },
      {
        "id": 4596,
        "field_group": {
          "id": 2349,
          "name": "UserDecisionInfo",
          "desc": "用户的决策信息：点击、消费、是否收藏",
          "ctime": 1729063127,
          "mtime": 1729063127
        },
        "extra": [
          {
            "id": 2263,
            "name": "product_order_info",
            "desc": "商品关联的订单信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2250,
            "name": "product_user_decision",
            "desc": "商品用户决策信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          },
          {
            "id": 2262,
            "name": "product_user_info",
            "desc": "商品关联的用户(点赞、购买)信息",
            "ctime": 1729063125,
            "mtime": 1729063125
          }
        ],
        "ctime": 1729063136,
        "mtime": 1729063136
      }
    ],
    "filters": [
      {
        "id": 1880,
        "name": "FilterByDiscountRevert",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 10,
        "psm": "data.life.commerce"
      },
      {
        "id": 1882,
        "name": "FilterByNotIndependentSale",
        "desc": "根据非独立销售信息过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 50000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1883,
        "name": "FilterByOutOfDate",
        "desc": "根据过期信息过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 47000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1886,
        "name": "FilterBySaleChannel",
        "desc": "根据销售渠道过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 56000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1887,
        "name": "FilterBySpuStatus",
        "desc": "根据SPU状态过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 43000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1888,
        "name": "FilterByLibraExperiment",
        "desc": "根据实验信息过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 54000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1889,
        "name": "FilterByMicroAppInvisibility",
        "desc": "根据微应用不可见性过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 45000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1891,
        "name": "FilterByProductStatus",
        "desc": "根据商品状态过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 46000,
        "psm": "data.life.commerce"
      },
      {
        "id": 1892,
        "name": "FilterByProductStock",
        "desc": "根据商品库存过滤",
        "plugin_id": 36,
        "plugin_name": "Marketing",
        "task_type": 3,
        "exec_type": 1,
        "ctime": 1729063129,
        "mtime": 1729611223,
        "priority": 71000,
        "psm": "data.life.commerce"
      }
    ],
    "ctime": 11,
    "mtime": 1729611223,
    "psm": "data.life.commerce"
  }
}
`
)
