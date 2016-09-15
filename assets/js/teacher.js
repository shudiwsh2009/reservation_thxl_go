var width=$(window).width();
var height=$(window).height();
var teacher;
var reservations;
var firstCategory;
var secondCategory;

function viewReservations() {
  getFeedbackCategories();
  $.getJSON('/teacher/reservation/view', function(json, textStatus) {
    if (json.state === 'SUCCESS') {
      console.log(json);
      reservations = json.reservations;
      teacher = json.teacher_info;
      refreshDataTable(reservations);
      optimize();
    } else {
      alert(json.message);
    }
  });
}

function getFeedbackCategories() {
  $.ajax({
    type: 'GET',
    async: false,
    url: '/category/feedback',
    dateType: 'json',
    success: function(data) {
      if (data.state === 'SUCCESS') {
        firstCategory = data.first_category;
        secondCategory = data.second_category;
      }
    }
  });
}

function refreshDataTable(reservations) {
  $('#page_maintable')[0].innerHTML = '\
    <div class="table_col" id="col_time">\
      <div class="table_head table_cell">时间</div>\
    </div>\
    <div class="table_col" id="col_teacher_fullname">\
      <div class="table_head table_cell">咨询师</div>\
    </div>\
    <div class="table_col" id="col_teacher_mobile">\
      <div class="table_head table_cell">咨询师手机</div>\
    </div>\
    <div class="table_col" id="col_status">\
      <div class="table_head table_cell">状态</div>\
    </div>\
    <div class="table_col" id="col_student">\
      <div class="table_head table_cell">学生</div>\
    </div>\
    <div class="clearfix"></div>\
  ';

  for (var i = 0; i < reservations.length; ++i) {
    $('#col_time').append('<div class="table_cell" id="cell_time_' + i + '" onclick="editReservation('
      + i + ')">' + reservations[i].start_time.split(' ')[0].substr(2) + '<br>'
      + reservations[i].start_time.split(' ')[1] + '-' + reservations[i].end_time.split(' ')[1] + '</div>');
    $('#col_teacher_fullname').append('<div class="table_cell" id="cell_teacher_fullname_'
      + i + '">' + reservations[i].teacher_fullname + '</div>');
    $('#col_teacher_mobile').append('<div class="table_cell" id="cell_teacher_mobile_'
      + i + '">' + reservations[i].teacher_mobile + '</div>');
    if (reservations[i].status === 'AVAILABLE') {
      $('#col_status').append('<div class="table_cell" id="cell_status_' + i + '">未预约</div>');
      $('#col_student').append('<div class="table_cell" id="cell_student_' + i + '">'
        + '<button type="button" id="cell_student_view_' + i + '" disabled="true" style="padding: 2px 2px">查看'
        + '</button></div>');
    } else if (reservations[i].status === 'RESERVATED') {
      $('#col_status').append('<div class="table_cell" id="cell_status_' + i + '">已预约</div>');
      $('#col_student').append('<div class="table_cell" id="cell_student_' + i + '">'
        + '<button type="button" id="cell_student_view_' + i + '" onclick="getStudent(' + i + ');" style="padding: 2px 2px">查看'
        + '</button></div>');
    } else if (reservations[i].status === 'FEEDBACK') {
      $('#col_status').append('<div class="table_cell" id="cell_status_' + i + '">'
        + '<button type="button" id="cell_status_feedback_' + i + '" onclick="getFeedback(' + i + ');" style="padding: 2px 2px">'
        + '反馈</button></div>');
      $('#col_student').append('<div class="table_cell" id="cell_student_' + i + '">'
        + '<button type="button" id="cell_student_view_' + i + '" onclick="getStudent(' + i + ');" style="padding: 2px 2px">查看'
        + '</button></div>');
    }
  }
}

function optimize(t) {
  $('#col_time').width(width * 0.25);
  $('#col_teacher_fullname').width(width * 0.2);
  $('#col_teacher_mobile').width(width * 0.25);
  $('#col_status').width(width * 0.13);
  $('#col_student').width(width * 0.13);
  $('#col_time').css('margin-left', width * 0.01 + 'px');
  for (var i = 0; i < reservations.length; ++i) {
    var maxHeight = Math.max(
        $('#cell_time_' + i).height(),
        $('#cell_teacher_fullname_' + i).height(),
        $('#cell_teacher_mobile_' + i).height(),
        $('#cell_status_' + i).height(),
        $('#cell_student_' + i).height()
      );
    $('#cell_time_' + i).height(maxHeight);
    $('#cell_teacher_fullname_' + i).height(maxHeight);
    $('#cell_teacher_mobile_' + i).height(maxHeight);
    $('#cell_status_' + i).height(maxHeight);
    $('#cell_student_' + i).height(maxHeight);

    if (i % 2 == 1) {
      $('#cell_time_' + i).css('background-color', 'white');
      $('#cell_teacher_fullname_' + i).css('background-color', 'white');
      $('#cell_teacher_mobile_' + i).css('background-color', 'white');
      $('#cell_status_' + i).css('background-color', 'white');
      $('#cell_student_' + i).css('background-color', 'white');
    } else {
      $('#cell_time_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_fullname_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_mobile_' + i).css('background-color', '#f3f3ff');
      $('#cell_status_' + i).css('background-color', '#f3f3ff');
      $('#cell_student_' + i).css('background-color', '#f3f3ff');
    }

    if (reservations[i].student_crisis_level && reservations[i].student_crisis_level !== '0') {
      //$('#cell_student_' + i).css('background-color', 'rgba(255, 0, 0, ' + parseInt(reservations[i].student_crisis_level) / 5 +')');
      $('#cell_student_view_' + i).css('background-color', 'rgba(255, 0, 0, ' + parseInt(reservations[i].student_crisis_level) / 1 +')');
    }
  }
  $(t).css('left', (width - $(t).width()) / 2 - 11 + 'px');
  $(t).css('top', (height - $(t).height()) / 2 - 11 + 'px');
}

function getFeedback(index) {
  $.post('/teacher/reservation/feedback/get', {
    reservation_id: reservations[index].reservation_id,
    source_id: reservations[index].source_id,
  }, function(data, textStatus, xhr) {
    if (data.state === 'SUCCESS') {
      showFeedback(index, data.feedback);
    } else {
      alert(data.message);
    }
  });
}

function showFeedback(index, feedback) {
  $('body').append('\
    <div class="pop_window" id="feedback_table_' + index + '" style="text-align: left; width: 90%; height: 70%; overflow:auto;">\
    咨询师反馈表<br>\
    评估分类：<br>\
    <select id="category_first_' + index + '" onchange="showSecondCategory(' + index + ')"><option value="">请选择</option></select><br>\
    <select id="category_second_' + index + '"></select><br>\
    出席人员：<br>\
    <input id="participant_student_' + index + '" type="checkbox">学生</input><input id="participant_parents_' + index + '" type="checkbox">家长</input>\
    <input id="participant_teacher_' + index + '" type="checkbox">教师</input><input id="participant_instructor_' + index + '" type="checkbox">辅导员</input>\
    <input id="participant_other_' + index + '" type="checkbox">其他</input><br>\
    重点明细：<select id="emphasis_'+ index + '"><option value="0">否</option><option value="1">是</option></select><br>\
    <div id="div_emphasis_' + index + '" style="display: none">\
      <b>严重程度：</b>\
      <input id="severity_' + index + '_0" type="checkbox">缓考</input>\
      <input id="severity_' + index + '_1" type="checkbox">休学复学</input>\
      <input id="severity_' + index + '_2" type="checkbox">家属陪读</input>\
      <input id="severity_' + index + '_3" type="checkbox">家属不知情</input>\
      <input id="severity_' + index + '_4" type="checkbox">任何其他需要知会院系关注的原因</input>\
      <br>\
      <b>疑似或明确的医疗诊断：</b>\
      <input id="medical_diagnosis_' + index + '_0" type="checkbox">服药</input>\
      <input id="medical_diagnosis_' + index + '_1" type="checkbox">精神分裂</input>\
      <input id="medical_diagnosis_' + index + '_2" type="checkbox">双相情感障碍</input>\
      <input id="medical_diagnosis_' + index + '_3" type="checkbox">焦虑症（状态）</input>\
      <br>　　　　　\
      <input id="medical_diagnosis_' + index + '_4" type="checkbox">抑郁症（状态）</input>\
      <input id="medical_diagnosis_' + index + '_5" type="checkbox">强迫症</input>\
      <input id="medical_diagnosis_' + index + '_6" type="checkbox">进食障碍</input>\
      <input id="medical_diagnosis_' + index + '_7" type="checkbox">失眠</input>\
      <input id="medical_diagnosis_' + index + '_8" type="checkbox">其他精神症状</input>\
      <input id="medical_diagnosis_' + index + '_9" type="checkbox">躯体疾病</input>\
      <input id="medical_diagnosis_' + index + '_10" type="checkbox">不遵医嘱</input>\
      <b>危急情况：</b>\
      <input id="crisis_' + index + '_0" type="checkbox">自伤</input>\
      <input id="crisis_' + index + '_1" type="checkbox">伤害他人</input>\
      <input id="crisis_' + index + '_2" type="checkbox">自杀念头</input>\
      <input id="crisis_' + index + '_3" type="checkbox">自杀未遂</input>\
    </div>\
    咨询记录：<br>\
    <textarea id="record_' + index + '" style="width: 100%; height:80px"></textarea><br>\
    是否危机个案：<select id="crisis_level_'+ index + '"><option value="0">否</option><option value="3">三星</option><option value="4">四星</option><option value="5">五星</option></select><br>\
    <button type="button" onclick="submitFeedback(' + index + ');">提交</button>\
    <button type="button" onclick="$(\'#feedback_table_' + index + '\').remove();">取消</button>\
    </div>\
  ');
  $(function() {
    showFirstCategory(index);
    if (feedback.category.length > 0) {
      $('#category_first_' + index).val(feedback.category.charAt(0));
      $('#category_first_' + index).change();
      $('#category_second_' + index).val(feedback.category);
    }
    if (feedback.participants.length > 0) {
      $('#participant_student_' + index).first().attr('checked', feedback.participants[0] > 0);
      $('#participant_parents_' + index).first().attr('checked', feedback.participants[1] > 0);
      $('#participant_teacher_' + index).first().attr('checked', feedback.participants[2] > 0);
      $('#participant_instructor_' + index).first().attr('checked', feedback.participants[3] > 0);
      $('#participant_other_' + index).first().attr('checked', feedback.participants[4] > 0);
    }
    var i = 1;
    for (i = 0; i < 5; i++) {
      $('#severity_' + index + '_' + i).first().attr('checked', feedback.severity[i] > 0);
    }
    for (i = 0; i < 11; i++) {
      $('#medical_diagnosis_' + index + '_' + i).first().attr('checked', feedback.medical_diagnosis[i] > 0);
    }
    for (i = 0; i < 4; i++) {
      $('#crisis_' + index + "_" + i).first().attr('checked', feedback.crisis[i] > 0);
    }
    $('#record_' + index).val(feedback.record);
    $('#emphasis_' + index).change(function() {
      if ($('#emphasis_' + index).val() === "0") {
        $('#div_emphasis_' + index).hide();
      } else {
        $('#div_emphasis_' + index).show();
      }
    });
    $('#emphasis_' + index).val(feedback.emphasis);
    $('#emphasis_' + index).change();
    $('#crisis_level_' + index).val(feedback.crisis_level);
    optimize('#feedback_table_' + index);
  });
}

function showFirstCategory(index) {
  for (var name in firstCategory) {
    if (firstCategory.hasOwnProperty(name)) {
      $('#category_first_' + index).append($("<option>", {
        value: name,
        text: firstCategory[name],
      }));
    }
  }
}

function showSecondCategory(index) {
  var first = $('#category_first_' + index).val();
  $('#category_second_' + index).find("option").remove().end().append('<option value="">请选择</option>').val('');
  if ($('#category_first_' + index).selectedIndex === 0) {
    return;
  }
  if (secondCategory.hasOwnProperty(first)) {
    for (var name in secondCategory[first]) {
      if (secondCategory[first].hasOwnProperty(name)) {
        var option = new Option(name, secondCategory[first][name]);
        $('#category_second_' + index).append($("<option>", {
          value: name,
          text: secondCategory[first][name],
        }));
      }
    }
  }
}

function submitFeedback(index) {
  var participants = [];
  participants.push($('#participant_student_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_parents_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_teacher_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_instructor_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_other_' + index).first().is(':checked') ? 1 : 0);
  var isEmphasis = $('#emphasis_' + index).val() !== "0"
  var i = 1;
  var severity = [];
  for (i = 0; i < 5; i++) {
    severity.push(isEmphasis ? ($('#severity_' + index + '_' + i).first().is(':checked') ? 1 : 0) : 0);
  }
  var medicalDiagnosis = [];
  for (i = 0; i < 11; i++) {
    medicalDiagnosis.push(isEmphasis ? ($('#medical_diagnosis_' + index + '_' + i).first().is(':checked') ? 1 : 0) : 0);
  }
  var crisis = [];
  for (i = 0; i < 4; i++) {
    crisis.push(isEmphasis ? ($('#crisis_' + index + "_" + i).first().is(':checked') ? 1 : 0) : 0);
  }
  var payload = {
    reservation_id: reservations[index].reservation_id,
    category: $('#category_second_' + index).val(),
    participants: participants,
    emphasis: $('#emphasis_' + index).val(),
    severity: severity,
    medical_diagnosis: medicalDiagnosis,
    crisis: crisis,
    record: $('#record_' + index).val(),
    crisis_level: $('#crisis_level_' + index).val(),
  };
  $.ajax({
    url: '/teacher/reservation/feedback/submit',
    type: 'POST',
    dataType: 'json',
    data: payload,
    traditional: true,
  })
  .done(function(data) {
    if (data.state === 'SUCCESS') {
      successFeedback(index);
      viewReservations();
    } else {
      alert(data.message);
    }
  });
}

function successFeedback(index) {
  $('#feedback_table_' + index).remove();
  $('body').append('\
    <div id="pop_success_feedback" class="pop_window" style="width: 50%;">\
      您已成功提交反馈！<br>\
      <button type="button" onclick="$(\'#pop_success_feedback\').remove();">确定</button>\
    </div>\
  ');
  optimize('#pop_success_feedback');
}

function queryStudent() {
  $.post('/teacher/student/query', {
    student_username: $('#query_student').val()
  }, function(data, textStatus, xhr) {
    if (data.state === 'SUCCESS') {
      showStudent(data.student_info, data.reservations);
    } else {
      alert(data.message);
    }
  });
}
function getStudent(index) {
  $.post('/teacher/student/get', {
    student_id: reservations[index].student_id
  }, function(data, textStatus, xhr) {
    if (data.state === 'SUCCESS') {
      showStudent(data.student_info, data.reservations);
    } else {
      alert(data.message);
    }
  });
}

function showStudent(student, reservations) {
  $('body').append('\
    <div id="pop_show_student_' + student.student_id + '" class="pop_window" style="text-align: left; width: 90%; height: 60%; overflow: auto">\
      学号：' + student.student_username + '<br>\
      姓名：' + student.student_fullname + '<br>\
      性别：' + student.student_gender + '<br>\
      出生日期：' + student.student_birthday + '<br>\
      系别：' + student.student_school + '<br>\
      年级：' + student.student_grade + '<br>\
      现住址：' + student.student_current_address + '<br>\
      家庭住址：' + student.student_family_address + '<br>\
      联系电话：' + student.student_mobile + '<br>\
      Email：' + student.student_email + '<br>\
      咨询经历：' + (student.student_experience_time ? '时间：' + student.student_experience_time + ' 地点：' + student.student_experience_location + ' 咨询师：' + student.student_experience_teacher : '无') + '<br>\
      父亲年龄：' + student.student_father_age + ' 职业：' + student.student_father_job + ' 学历：' + student.student_father_edu + '<br>\
      母亲年龄：' + student.student_mother_age + ' 职业：' + student.student_mother_job + ' 学历：' + student.student_mother_edu + '<br>\
      父母婚姻状况：' + student.student_parent_marriage + '<br>\
      近三个月里发生的有重大意义的事：' + student.student_significant + '<br>\
      需要接受帮助的主要问题：' + student.student_problem + '<br>\
      档案分类：' + student.student_archive_category + ' 档案编号：' + student.student_archive_number + '<br>\
      是否危机个案：<span id="crisis_level_'+ student.student_id + '"></span><br>\
      已绑定的咨询师：<span id="binded_teacher_username">' + student.student_binded_teacher_username + '</span>&nbsp;\
        <span id="binded_teacher_fullname">' + student.student_binded_teacher_fullname + '</span><br>\
      <br>\
      <button type="button" onclick="$(\'#pop_show_student_' + student.student_id + '\').remove();">关闭</button>\
      <div id="student_reservations_' + student.student_id + '" style="width: 100%">\
      </div>\
    </div>\
  ');
  for (var i = 0; i < reservations.length; i++) {
    $('#student_reservations_' + student.student_id).append('\
      <div class="has_children" style="background: ' + (reservations[i].status === 'FEEDBACK' ? '#555' : '#F00') + '">\
        <span>' + reservations[i].start_time + ' 至 ' + reservations[i].end_time + '  ' + reservations[i].teacher_fullname + '</span>\
        <p class="children">学生反馈：' + reservations[i].student_feedback.scores + '</p>\
        <p class="children">评估分类：' + reservations[i].teacher_feedback.category + '</p>\
        <p class="children">出席人员：' + reservations[i].teacher_feedback.participants + '</p>\
        <p class="children">重点明细：' + (reservations[i].teacher_feedback.emphasis === '0' ? '否' : '是') + '</p>'
        + (reservations[i].teacher_feedback.severity === '' ? '' : '<p class="children">严重程度：' + reservations[i].teacher_feedback.severity + '</p>')
        + (reservations[i].teacher_feedback.medical_diagnosis === '' ? '' : '<p class="children">疑似或明确的医疗诊断：' + reservations[i].teacher_feedback.medical_diagnosis + '</p>')
        + (reservations[i].teacher_feedback.crisis === '' ? '' : '<p class="children">危急情况：' + reservations[i].teacher_feedback.crisis + '</p>')
        + '<p class="children">咨询记录：' + reservations[i].teacher_feedback.record + '</p>\
      </div>\
    ');
  }
  $(function() {
    $('.has_children').click(function() {
      $(this).addClass('highlight').children('p').show().end()
          .siblings().removeClass('highlight').children('p').hide();
    });
    $('#crisis_level_' + student.student_id).text(student.student_crisis_level === 0 ? '否' : '是');
  });
  optimize('#pop_show_student_' + student.student_id);
}
