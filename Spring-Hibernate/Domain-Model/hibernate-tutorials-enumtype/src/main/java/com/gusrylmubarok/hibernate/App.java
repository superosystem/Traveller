package com.gusrylmubarok.hibernate;

import com.gusrylmubarok.hibernate.dao.ProjectDao;
import com.gusrylmubarok.hibernate.entity.Project;
import com.gusrylmubarok.hibernate.entity.ProjectStatus;

import java.util.List;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main( String[] args )
    {
        ProjectDao projectDao = new ProjectDao();
        Project project = new Project();
        project.setProjectName("OOP");
        project.setProjectStatus(ProjectStatus.INPROGRESS);
        projectDao.saveProject(project);

        List<Project> projects = projectDao.getProjects();
        projects.forEach(s -> {
            System.out.println(s.getProjectName());
            System.out.println(s.getProjectStatus());
        });
    }
}
