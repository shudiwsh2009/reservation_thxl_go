import 'weui';
import { AlertDialog, LoadingHud } from '#coms/Huds';
import { Application, User } from '#models/Models';
import { Panel, PanelHeader } from 'react-weui';
import PageBottom from '#coms/PageBottom';
import PropTypes from 'prop-types';
import React from 'react';
import TeacherFeedbackForm from '#forms/TeacherFeedbackForm';
import queryString from 'query-string';

export default class TeacherFeedbackPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservation: null,
      feedback: null,
      student: null,
      crisisLevel: 0,
    };
    this.handleCancel = this.handleCancel.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.showAlert = this.showAlert.bind(this);
  }

  componentDidMount() {
    this.loading.show('正在加载中');
    const parsedQuery = queryString.parse(this.props.history.location.search);
    const reservationId = parsedQuery.reservation_id;
    User.updateSession(() => {
      Application.getFeedbackByTeacher(reservationId, (data) => {
        setTimeout(() => {
          this.setState({
            feedback: data.feedback,
            student: data.student,
            reservation: data.reservation,
          }, () => {
            this.loading.hide();
          });
        }, 500);
      }, (error) => {
        this.loading.hide();
        this.alert.show('', error, '好的', () => {
          this.props.history.push('/reservation');
        });
      });
    }, () => {
      this.props.history.push('/login');
    });
  }

  handleCancel() {
    this.props.history.goBack();
  }

  handleSubmit(payload) {
    this.loading.show('正在加载中');
    setTimeout(() => {
      Application.submitFeedbackByTeacher(payload, () => {
        this.loading.hide();
        this.alert.show('提交成功', '您已成功提交反馈', '好的', () => {
          this.props.history.push('/reservation');
        });
      }, (error) => {
        this.loading.hide();
        this.alert.show('提交失败', error, '好的');
      });
    }, 500);
  }

  showAlert(title, msg, label, click) {
    this.alert.show(title, msg, label, click);
  }

  render() {
    return (
      <div>
        <Panel>
          <PanelHeader style={{ fontSize: '18px' }}>{this.state.student && `${this.state.student.fullname}同学的`}咨询师反馈表</PanelHeader>
          <TeacherFeedbackForm
            reservation={this.state.reservation}
            feedback={this.state.feedback}
            handleSubmit={this.handleSubmit}
            handleCancel={this.handleCancel}
            showAlert={this.showAlert}
          />
          <LoadingHud ref={(loading) => { this.loading = loading; }} />
          <AlertDialog ref={(alert) => { this.alert = alert; }} />
          <PageBottom
            styles={{ color: '#999999', textAlign: 'center', backgroundColor: 'white', fontSize: '14px' }}
            contents={['清华大学学生心理发展指导中心', '联系方式：010-62782007']}
            height="55px"
          />
        </Panel>
      </div>
    );
  }
}

TeacherFeedbackPage.propTypes = {
  history: PropTypes.object.isRequired,
};
