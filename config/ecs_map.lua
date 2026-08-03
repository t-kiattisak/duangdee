-- Maps parsed zerolog JSON fields to Elastic Common Schema (ECS) field names.

local function trim_container_prefix(name)
    if name:sub(1, 1) == "/" then
        return name:sub(2)
    end
    return name
end

local function encode_query(params)
    if type(params) ~= "table" then
        return tostring(params)
    end
    local parts = {}
    for k, v in pairs(params) do
        table.insert(parts, k .. "=" .. tostring(v))
    end
    table.sort(parts)
    return table.concat(parts, "&")
end

local function container_fields(record, target)
    if record["container_name"] then
        target["container.name"] = trim_container_prefix(record["container_name"])
    end
    if record["container_id"] then
        target["container.id"] = record["container_id"]
    end
    if record["source"] then
        target["log.origin.file.name"] = record["source"]
    end
end

function ecs_map(tag, timestamp, record)
    local ecs = {}

    container_fields(record, ecs)

    -- Structured JSON log emitted by zerolog
    if record["level"] ~= nil then
        if record["time"] then
            ecs["@timestamp"] = record["time"]
        end

        ecs["log.level"] = record["level"]
        ecs["message"] = record["message"] or ""
        ecs["event.module"] = "duangdee"
        ecs["event.kind"] = "event"
        ecs["event.category"] = {"web"}
        ecs["event.type"] = {"access"}

        if record["service"] then
            ecs["service.name"] = record["service"]
            ecs["event.dataset"] = "duangdee." .. record["service"]
        end
        if record["trace_id"] then
            ecs["trace.id"] = record["trace_id"]
        end
        if record["http_method"] then
            ecs["http.request.method"] = record["http_method"]
        end
        if record["http_path"] then
            ecs["url.path"] = record["http_path"]
        end
        if record["http_status"] then
            ecs["http.response.status_code"] = record["http_status"]
        end
        if record["client_ip"] then
            ecs["client.ip"] = record["client_ip"]
        end
        if record["user_agent"] then
            ecs["user_agent.original"] = record["user_agent"]
        end
        if record["request_headers"] then
            ecs["http.request.headers"] = record["request_headers"]
        end
        if record["request_payload"] then
            ecs["http.request.body.content"] = record["request_payload"]
        end
        if record["response_payload"] then
            ecs["http.response.body.content"] = record["response_payload"]
        end
        if record["query_params"] then
            ecs["url.query"] = encode_query(record["query_params"])
        end

        return 2, timestamp, ecs
    end

    -- Plain-text fallback for unstructured container logs
    if record["log"] then
        ecs["message"] = record["log"]
        ecs["log.level"] = "info"
        ecs["event.module"] = "duangdee"

        if ecs["container.name"] then
            ecs["service.name"] = ecs["container.name"]
            ecs["event.dataset"] = "duangdee." .. ecs["container.name"]
        end

        return 2, timestamp, ecs
    end

    return 1, timestamp, record
end
