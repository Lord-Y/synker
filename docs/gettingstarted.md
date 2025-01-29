# Getting started

`Synker` currently support only `json` and `yaml` files.

## Schemas

The schemas block in the config file must have 4 blocks:
- `topic` block permit to create the topic with the required configuration
- `changeFeed` block permit create the `change feed` on `CockroachDB`
- `elasticsearch` hold the `elasticsearch` config with the required mapping. if `elasticsearch.index.create` is set to `true`, the index will be created otherwise no.
- `sql` block hold which query type need to be applied on between `CockroachDB` and `elasticsearch`.
- `queryType` hold the requirements needed to apply when a message is received in the `topic`

## Query type `none`

With the query type `none`, no SQL queries will be performed on the database but it contains the immutable colums that permit `Synker` to create a uniq data in `elasticsearch` and update this data if necessary

Here is an example:
```yaml
schemas:
- name: promo_codes
  topic:
    name: movr.public.promo_codes
    numPartitions: 1
    replicationFactor: 3
    config:
    - key: max.message.bytes
      value: "128000"
  changeFeed:
    fullTableName: movr.public.promo_codes
    options:
    - full_table_name
    - on_error = 'pause'
    - protect_data_from_gc_on_pause
    - updated
  sql:
    queryType:
      none:
        immutableColumns:
        - name: code
          type: varchar
  elasticsearch:
    index:
      name: promo_codes
      create: true
    mapping:
      settings:
        index.requests.cache.enable: true
        index:
          number_of_shards: 3
          number_of_replicas: 2
        analysis:
          normalizer:
            lowerasciinormalizer:
              type: custom
              filter:
              - lowercase
              - asciifolding
```

## SQL Parser

We won't implement an `sql parser` because it is an heavy things to do.
All `SQL queries` must be:
- `SELECT table_name.column_a,table_name.column_b`
- `SELECT table_name.column_a,table_name.column_b FROM table_name WHERE (table_name.column_c = 1)`
- must not contain any aliases on tables: `SELECT a.column_a,a.column_b FROM table_name a` because:
  - we need to update the query filter with data coming from `kafka message` and having some aliases can lead to weird queries that won't be easy to debug.

We also won't implement a `golang sql query structure` to build the query because `SELECT` queries can contain a lot of logics and it will be difficult to support all of them.

## Query type `advanced`

With the query type `advanced` you can perform complex SQL queries after receiving a message.

Here is an example:
```yaml
- name: rides
  topic:
    name: movr.public.rides
    numPartitions: 1
    replicationFactor: 3
    config:
    - key: max.message.bytes
      value: "128000"
  changeFeed:
    fullTableName: movr.public.rides
    options:
    - full_table_name
    - on_error = 'pause'
    - protect_data_from_gc_on_pause
    - updated
  sql:
    queryType:
      advanced:
        immutableColumns:
        - name: id
          type: uuid
        query: "SELECT rides.id,rides.city,rides.vehicle_city,rides.rider_id,rides.vehicle_id,rides.start_address,rides.end_address,rides.start_time,rides.end_time,rides.revenue,vehicles.type FROM rides LEFT JOIN vehicles ON vehicles.id = rides.vehicle_id"
  elasticsearch:
    index:
      name: rides
      create: true
    mapping:
      settings:
        index.requests.cache.enable: true
        index:
          number_of_shards: 3
          number_of_replicas: 2
        analysis:
          normalizer:
            lowerasciinormalizer:
              type: custom
              filter:
              - lowercase
              - asciifolding
      mappings:
        dynamic_templates:
        - string_as_keyword:
            match_mapping_type: string
            match: '*'
            mapping:
              type: keyword
              normalizer: lowerasciinormalizer
        properties:
          id:
            type: keyword
          city:
            normalizer: lowerasciinormalizer
            type: keyword
          vehicle_city:
            normalizer: lowerasciinormalizer
            type: keyword
          rider_id:
            type: keyword
          vehicle_id:
            type: keyword
          start_address:
            normalizer: lowerasciinormalizer
            type: keyword
          end_address:
            normalizer: lowerasciinormalizer
            type: keyword
          start_time:
            type: date
          end_time:
            type: date
          revenue:
            type: float
          type:
            normalizer: lowerasciinormalizer
            type: keyword
```

If you check the SQL query `SELECT rides.id,rides.city,rides.vehicle_city,rides.rider_id,rides.vehicle_id,rides.start_address,rides.end_address,rides.start_time,rides.end_time,rides.revenue,vehicles.type FROM rides LEFT JOIN vehicles ON vehicles.id = rides.vehicle_id`, you can see that it's a join query between tables `rides` and `vehicules`.
Fortunately when updates will be done on `vehicules`, we will also received a kafka message that will be then processed by `synker`.

## Examples

More examples are present in the `processing/examples/schemas` folder.
