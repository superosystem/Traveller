package com.gusrylmubarok.hibernate.dao;

import com.gusrylmubarok.hibernate.entity.Project;
import com.gusrylmubarok.hibernate.util.HibernateUtil;
import org.hibernate.Session;
import org.hibernate.Transaction;

import java.util.List;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class ProjectDao {

    public void saveProject(Project project) {
        try (Session session = HibernateUtil.getSessionFactory().openSession()) {
            // Start
            Transaction transaction = session.beginTransaction();
            // save project
            session.save(project);
            // commit transaction
            transaction.commit();
        }
    }

    public List<Project> getProjects() {
        try (Session session = HibernateUtil.getSessionFactory().openSession()) {
            return session.createQuery("from Project", Project.class).list();
        }
    }

}
