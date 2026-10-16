import dayjs from 'dayjs';

const datetime = {
  format(dateString, format = 'DD/MM/YYYY') {
    console.log('dateString:', dateString);
    return dayjs(dateString).format(format);
  },
};

export default datetime;
