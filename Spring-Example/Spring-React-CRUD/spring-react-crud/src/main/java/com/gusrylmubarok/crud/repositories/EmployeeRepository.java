package com.gusrylmubarok.crud.repositories;

import com.gusrylmubarok.crud.entities.Employee;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Repository
public interface EmployeeRepository extends JpaRepository<Employee, Long> {
}
