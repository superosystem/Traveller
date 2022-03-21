package com.gusrylmubarok.hibernate;

import com.gusrylmubarok.hibernate.dao.StudentDao;
import com.gusrylmubarok.hibernate.entity.Student;

import java.util.List;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main(String[] args) {
        StudentDao studentDao = new StudentDao();
        Student student = new Student("Gusryl", "Mubarok", "gusrylmubarok@gmail.com");
        studentDao.saveStudent(student);

        List< Student > students = studentDao.getStudents();
        students.forEach(s -> System.out.println(s.getFirstName()));
    }
}
