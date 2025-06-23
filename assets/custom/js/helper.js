/**
 * 渲染客户端列表
 * @param {Array} data - 包含客户端数据的数组
 */
function renderClientList(data) {
    $.each(data, function(index, obj) {
        // 渲染Telegram按钮
        let telegramButton = ''
        if (obj.Client.telegram_userid) {
            telegramButton =    `<div class="btn-group">      
                                    <button type="button" class="btn btn-outline-primary btn-sm" data-toggle="modal"
                                        data-target="#modal_telegram_client" data-clientid="${obj.Client.id}"
                                        data-clientname="${obj.Client.name}">Telegram</button>
                                </div>`
        }

        let telegramHtml = "";
        if (obj.Client.telegram_userid && obj.Client.telegram_userid.length > 0) {
            telegramHtml = `<span class="info-box-text" style="display: none"><i class="fas fa-tguserid"></i>${obj.Client.telegram_userid}</span>`
        }

        // 渲染客户端状态的CSS标签样式
        let clientStatusHtml = '>'
        if (obj.Client.enabled) {
            clientStatusHtml = `style="visibility: hidden;">`
        }

        // 渲染客户端已分配的IP地址
        let allocatedIpsHtml = "";
        $.each(obj.Client.allocated_ips, function(index, obj) {
            allocatedIpsHtml += `<small class="badge badge-secondary">${obj}</small>&nbsp;`;
        })

        // 渲染客户端允许的IP地址
        let allowedIpsHtml = "";
        $.each(obj.Client.allowed_ips, function(index, obj) {
            allowedIpsHtml += `<small class="badge badge-secondary">${obj}</small>&nbsp;`;
        })

        let subnetRangesString = "";
        if (obj.Client.subnet_ranges && obj.Client.subnet_ranges.length > 0) {
            subnetRangesString = obj.Client.subnet_ranges.join(',')
        }

        let additionalNotesHtml = "";
        if (obj.Client.additional_notes && obj.Client.additional_notes.length > 0) {
            additionalNotesHtml = `<span class="info-box-text" style="display: none"><i class="fas fa-additional_notes"></i>${obj.Client.additional_notes.toUpperCase()}</span>`
        }

        // 渲染客户端HTML内容
        let html = `<div class="col-sm-6 col-md-6 col-lg-4" id="client_${obj.Client.id}">
                        <div class="info-box">
                            <div class="overlay" id="paused_${obj.Client.id}"` + clientStatusHtml
                                + `<i class="paused-client fas fa-3x fa-play" onclick="resumeClient('${obj.Client.id}')"></i>
                            </div>
                            <div class="info-box-content" style="overflow: hidden">
                                <div class="btn-group">
                                    <a href="download?clientid=${obj.Client.id}" class="btn btn-outline-primary btn-sm">下载</a>
                                </div>
                                <div class="btn-group">      
                                    <button type="button" class="btn btn-outline-primary btn-sm" data-toggle="modal"
                                        data-target="#modal_qr_client" data-clientid="${obj.Client.id}"
                                        data-clientname="${obj.Client.name}" ${obj.QRCode != "" ? '' : ' disabled'}>二维码</button>
                                </div>
                                <div class="btn-group">      
                                    <button type="button" class="btn btn-outline-primary btn-sm" data-toggle="modal"
                                        data-target="#modal_email_client" data-clientid="${obj.Client.id}"
                                        data-clientname="${obj.Client.name}">发送邮件</button>
                                </div>
                                ${telegramButton}
                                <div class="btn-group">
                                    <button type="button" class="btn btn-outline-danger btn-sm">更多</button>
                                    <button type="button" class="btn btn-outline-danger btn-sm dropdown-toggle dropdown-icon" 
                                        data-toggle="dropdown">
                                    </button>
                                    <div class="dropdown-menu" role="menu">
                                        <a class="dropdown-item" href="#" data-toggle="modal"
                                        data-target="#modal_edit_client" data-clientid="${obj.Client.id}"
                                        data-clientname="${obj.Client.name}">编辑</a>
                                        <a class="dropdown-item" href="#" data-toggle="modal"
                                        data-target="#modal_pause_client" data-clientid="${obj.Client.id}"
                                        data-clientname="${obj.Client.name}">禁用</a>
                                        <a class="dropdown-item" href="#" data-toggle="modal"
                                        data-target="#modal_remove_client" data-clientid="${obj.Client.id}"
                                        data-clientname="${obj.Client.name}">删除</a>
                                    </div>
                                </div>
                                <hr>
                                <span class="info-box-text"><i class="fas fa-user"></i> ${obj.Client.name}</span>
                                <span class="info-box-text" style="display: none"><i class="fas fa-key"></i> ${obj.Client.public_key}</span>
                                <span class="info-box-text" style="display: none"><i class="fas fa-subnetrange"></i>${subnetRangesString}</span>
                                ${telegramHtml}
                                ${additionalNotesHtml}
                                <span class="info-box-text"><i class="fas fa-envelope"></i> ${obj.Client.email}</span>
                                <span class="info-box-text"><i class="fas fa-clock"></i>
                                    ${prettyDateTime(obj.Client.created_at)}</span>
                                <span class="info-box-text"><i class="fas fa-history"></i>
                                    ${prettyDateTime(obj.Client.updated_at)}</span>
                                <span class="info-box-text"><i class="fas fa-server" style="${obj.Client.use_server_dns ? "opacity: 1.0" : "opacity: 0.5"}"></i>
                                    ${obj.Client.use_server_dns ? '启用DNS' : '禁用DNS'}</span>
                                <span class="info-box-text"><i class="fas fa-file"></i>
                                    ${obj.Client.additional_notes}</span>
                                <span class="info-box-text"><strong>IP分配</strong></span>`
                                + allocatedIpsHtml
                                + `<span class="info-box-text"><strong>允许的IP</strong></span>`
                                + allowedIpsHtml
                            +`</div>
                        </div>
                    </div>`

        // 将客户端HTML元素添加到列表中
        $('#client-list').append(html);
    });
}

/**
 * 渲染用户列表
 * @param {Array} data - 包含用户数据的数组
 */
function renderUserList(data) {
    $.each(data, function(index, obj) {
        let clientStatusHtml = '>'

        // 渲染用户HTML内容
        let html = `<div class="col-sm-6 col-md-6 col-lg-4" id="user_${obj.username}">
                        <div class="info-box">
                            <div class="info-box-content">
                                <div class="btn-group">
                                     <button type="button" class="btn btn-outline-primary btn-sm" data-toggle="modal" data-target="#modal_edit_user" data-username="${obj.username}">编辑</button>
                                </div>
                                <div class="btn-group">
                                    <button type="button" class="btn btn-outline-danger btn-sm" data-toggle="modal"
                                        data-target="#modal_remove_user" data-username="${obj.username}">删除</button>
                                </div>
                                <hr>
                                <span class="info-box-text"><i class="fas fa-user"></i> ${obj.username}</span>
                                <span class="info-box-text"><i class="fas fa-terminal"></i> ${obj.admin? '管理员' : '经理'}</span>
                                </div>
                        </div>
                    </div>`

        // 将用户HTML元素添加到列表中
        $('#users-list').append(html);
    });
}


/**
 * 格式化日期时间
 * @param {string} timeStr - 时间字符串
 * @returns {string} 格式化后的本地时间字符串
 */
function prettyDateTime(timeStr) {
    const dt = new Date(timeStr);
    const offsetMs = dt.getTimezoneOffset() * 60 * 1000;
    const dateLocal = new Date(dt.getTime() - offsetMs);
    return dateLocal.toISOString().slice(0, 19).replace(/-/g, "/").replace("T", " ");
}
