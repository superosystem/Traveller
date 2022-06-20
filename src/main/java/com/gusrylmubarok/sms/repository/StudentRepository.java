package com.gusrylmubarok.sms.repository;

import com.gusrylmubarok.sms.model.Student;
import org.springframework.data.jpa.repository.JpaRepository;

public interface StudentRepository extends JpaRepository<Student, Long> {
}
